package schedule

import (
	"fmt"
	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/register"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/redis"
)

func Init() {
	client := redis.ClientRedis

	locker, err := redislock.NewRedisLocker(client, redislock.WithTries(config.DefaultScheduleLockRedisRetry))
	if err != nil {
		panic(err)
	}

	s, err := gocron.NewScheduler(gocron.WithDistributedLocker(locker))
	if err != nil {
		panic(err)
	}

	for _, schedule := range register.Schedules {
		_, err = s.NewJob(
			schedule.JobDefinition,
			schedule.Task,
			schedule.Options...,
		)
		if err != nil {
			panic(err)
		}
	}

	s.Start()
	fmt.Println("Success init schedule/cron")
}
