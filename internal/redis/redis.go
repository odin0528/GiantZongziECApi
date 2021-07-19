package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("连接redis出错，错误信息：%v", err)
	}
	fmt.Println("成功连接redis")
}

func Set(key string, value string, times int) {
	rdb.Set(ctx, key, value, 10*time.Second)
}

func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}
