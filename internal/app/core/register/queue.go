package register

import "github.com/hibiken/asynq"

// Handlers is the extension point for application-specific background jobs.
var Handlers = map[string]asynq.HandlerFunc{}
