package register

import "github.com/go-co-op/gocron/v2"

type Schedule struct {
	JobDefinition gocron.JobDefinition
	Task          gocron.Task
	Options       []gocron.JobOption
}

// Schedules is intentionally empty in the base application.
var Schedules = []Schedule{}
