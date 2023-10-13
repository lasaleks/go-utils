package goutils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type User struct {
	Id          int64  `json:"id"`
	UserName    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Password    string `json:"password"`
	EMail       string `json:"email"`
	IsSuperUser bool   `json:"is_superuser"`
	IsStaff     bool   `json:"is_staff"`
	IsActive    bool   `json:"is_active"`
	DateJoined  string `json:"date_joined"`
	LastLogin   string `json:"last_login"`
	//"{\"id\": 8, \"password\": \"c4ca4238a0b923820dcc509a6f75849b\", \"last_login\": \"2020-11-25T06:48:12Z\", \"is_superuser\": false, \"username\": \"user\", \"first_name\": \"\", \"last_name\": \"\", \"email\": \"\", \"is_staff\": false, \"is_active\": true, \"date_joined\": \"2020-07-06T08:47:16Z\"}"
}

func GetDjangoUser(ctx context.Context, user_id int64, redis_cli *redis.Client) (user *User, err error) {
	json_data, err := redis_cli.HGet(ctx, "django.auth.user", fmt.Sprintf("%d", user_id)).Result()
	if err != nil {
		return nil, err
	}
	user = &User{}
	err = json.Unmarshal([]byte(json_data), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
