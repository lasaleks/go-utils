package django_sessions

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nlpodyssey/gopickle/pickle"
	"github.com/nlpodyssey/gopickle/types"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type keyContext int

const userIDKey keyContext = 1

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
			user_id, err = strconv.ParseInt(user_id_str, 10, 64)
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
		user_id := ctx.Value(userIDKey)
		ismem, err := rdb.SIsMember(ctx, fmt.Sprintf("auth.user.permission.%v", user_id), perm).Result()
		if err != nil || !ismem {
			http.Error(w, fmt.Sprintf("access denied %s", r.URL.Path), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIdInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	mdOut := metadata.Pairs(
		"user_id", fmt.Sprintf("%d", ctx.Value(userIDKey)),
	)
	callContext := metadata.NewOutgoingContext(ctx, mdOut)
	err := invoker(callContext, method, req, reply, cc, opts...)
	return err
}

/*
func AccessLogInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	md,_:=metadata.FromOutgoingContext(ctx)
	start:=time.Now()

	var traceId,userId,userRole string
	if len(md["authorization"])>0{
		tokenString:= md["authorization"][0]
		if tokenString!=""{
			err,token:=userService.CheckGetJWTToken(tokenString)
			if err!=nil{
				return err
			}
			userId=fmt.Sprintf("%s",token["UserID"])
			userRole=fmt.Sprintf("%s",token["UserRole"])
		}
	}
	//Присваиваю ID запроса
	traceId=fmt.Sprintf("%d",time.Now().UTC().UnixNano())

	callContext:=context.Background()
	mdOut:=metadata.Pairs(
		"trace-id",traceId,
		"user-id",userId,
		"user-role",userRole,
	)
	callContext=metadata.NewOutgoingContext(callContext,mdOut)

	err:=invoker(callContext,method,req,reply,cc, opts...)

	msg:=fmt.Sprintf("Call:%v, traceId: %v, userId: %v, userRole: %v, time: %v", method,traceId,userId,userRole,time.Since(start))
	app.AccesLog(msg)

	return err
}*/
