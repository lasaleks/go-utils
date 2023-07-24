package django_sessions

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nlpodyssey/gopickle/pickle"
	"github.com/nlpodyssey/gopickle/types"
	"github.com/redis/go-redis/v9"
)

var userIDKey = "user_id"

func SessionMiddleware(next http.Handler, rdb1 *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionid, err := r.Cookie("sessionid")
		if err != nil {
			http.Error(w, fmt.Sprintf("no auth at %s", r.URL.Path), http.StatusUnauthorized)
			return
		}
		KEY_PREFIX := "django.contrib.sessions.cache"
		version := 1

		session_key := fmt.Sprintf("%s:%d:%s%s", "", version, KEY_PREFIX, sessionid.Value)
		result, err := rdb1.Get(ctx, session_key).Result()
		if err == redis.Nil {
			http.Error(w, fmt.Sprintf("no auth at %s", r.URL.Path), http.StatusUnauthorized)
			return
		}
		py_i, err := pickle.Loads(result)
		if err != nil {
			http.Error(w, fmt.Sprintf("no auth at %s", r.URL.Path), http.StatusUnauthorized)
			return
		}
		var sess_data *types.Dict
		var ok bool

		if sess_data, ok = py_i.(*types.Dict); !ok {
			http.Error(w, fmt.Sprintf("no auth at %s", r.URL.Path), http.StatusUnauthorized)
			return
		}
		var user_id int64 = 0
		user_id_i, ok := sess_data.Get("_auth_user_id")
		if !ok {
			http.Error(w, fmt.Sprintf("no auth at %s", r.URL.Path), http.StatusUnauthorized)
			return
		}
		user_id_str, ok := user_id_i.(string)
		if ok {
			user_id, err = strconv.ParseInt(user_id_str, 10, 32)
			if err != nil {
				http.Error(w, fmt.Sprintf("no auth at %s", r.URL.Path), http.StatusUnauthorized)
				return
			}
		}

		ctx = context.WithValue(ctx, userIDKey, user_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ChkPermMiddleware(next http.Handler, rdb *redis.Client, perm string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ismem, err := rdb.SIsMember(ctx, fmt.Sprintf("auth.user.permission.%v", ctx.Value(userIDKey)), perm).Result()
		if err != nil || !ismem {
			http.Error(w, fmt.Sprintf("access denied %s", r.URL.Path), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
