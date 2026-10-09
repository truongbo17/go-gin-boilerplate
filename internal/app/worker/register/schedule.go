package register

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/schedule"
)

func exampleSchedules() []schedule.Job {
	return []schedule.Job{{
		JobDefinition: gocron.CronJob("*/15 * * * * *", true),
		Task:          gocron.NewTask(func() { fmt.Println("schedule") }),
		Options:       []gocron.JobOption{gocron.WithName("example:print-schedule")},
	}}
}
