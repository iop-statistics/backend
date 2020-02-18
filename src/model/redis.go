package model

import (
	"fmt"
	"github.com/Lyt99/iop-statistics/config"
	"github.com/go-redis/redis"
)

var (
	RedisClient *redis.Client
)

func init() {
	addr := fmt.Sprintf("%s:%s", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	err := RedisClient.Ping().Err()
	if err != nil {
		panic(err)
	}
}
