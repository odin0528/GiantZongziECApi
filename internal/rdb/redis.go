package rdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var ctx = context.Background()

func init() {

	RDB = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "",
		DB:       0,
		PoolSize: 1000,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("redis connect waring message:%v", err)
	} else {
		fmt.Println("redis connected")
	}
}

func Set(key string, value interface{}) {
	str, _ := json.Marshal(value)
	RDB.Set(ctx, key, string(str), 24*60*60*time.Second)
}

func Get(key string, res interface{}) error {
	result, err := RDB.Get(ctx, key).Result()
	json.Unmarshal([]byte(result), res)
	return err
}

func Del(key string) error {
	return RDB.Del(ctx, key).Err()
}
