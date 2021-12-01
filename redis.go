package main

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}
