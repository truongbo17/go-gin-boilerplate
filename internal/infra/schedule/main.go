package schedule

import (
	"fmt"

	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	redisclient "github.com/redis/go-redis/v9"
	"github.com/truongbo17/go-gin-boilerplate/config"
)

type Job struct {
	JobDefinition gocron.JobDefinition
	Task          gocron.Task
	Options       []gocron.JobOption
}

func StartWith(client *redisclient.Client, jobs []Job) (gocron.Scheduler, error) {
	locker, err := redislock.NewRedisLocker(client, redislock.WithTries(config.DefaultScheduleLockRedisRetry))
	if err != nil {
		return nil, fmt.Errorf("create schedule lock: %w", err)
	}

	s, err := gocron.NewScheduler(gocron.WithDistributedLocker(locker))
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	for _, schedule := range jobs {
		_, err = s.NewJob(
			schedule.JobDefinition,
			schedule.Task,
			schedule.Options...,
		)
		if err != nil {
			_ = s.Shutdown()
			return nil, fmt.Errorf("register schedule: %w", err)
		}
	}

	s.Start()
	fmt.Println("Success init schedule/cron")
	return s, nil
}
