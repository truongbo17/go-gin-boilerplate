package auth

import (
	"github.com/sirupsen/logrus"
	coregraphql "github.com/truongbo17/go-gin-boilerplate/internal/app/core/graphql"
)

type Resolver struct {
	Metadata coregraphql.Service
	Logger   *logrus.Logger
}
