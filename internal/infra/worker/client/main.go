package client

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/truongbo17/go-gin-boilerplate/config"
)

var WorkerClient *asynq.Client

func InitClient() {
	EnvConfig := config.EnvConfig
	configRedis := EnvConfig.Cache

	rdb := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", configRedis.RedisHost, configRedis.RedisPort),
		Username: configRedis.RedisUsername,
		Password: configRedis.RedisPassword,
	})

	WorkerClient = rdb

	fmt.Println("Success init client asynq queue.")
}
