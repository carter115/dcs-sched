package db

import (
	"dcs-sched/master/config"
	"github.com/go-redis/redis"
)

var RedisCli *redis.Client

func init() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.Db,
	})
}
