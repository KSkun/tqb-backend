package model

import (
	"github.com/go-redis/redis/v7"
	"github.com/KSkun/tqb-backend/config"
)

// 使用文档 https://redis.uptrace.dev/#executing-commands

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     config.C.Redis.Addr,
			Password: config.C.Redis.Password,
			DB:       config.C.Redis.DB,
		},
	)

	err := redisClient.Ping().Err()
	if err != nil {
		panic(err)
	}
}
