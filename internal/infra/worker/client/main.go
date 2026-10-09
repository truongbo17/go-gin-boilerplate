package client

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/truongbo17/go-gin-boilerplate/config"
)

func New(configRedis config.Cache) *asynq.Client {
	rdb := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", configRedis.RedisHost, configRedis.RedisPort),
		Username: configRedis.RedisUsername,
		Password: configRedis.RedisPassword,
	})

	return rdb
}
