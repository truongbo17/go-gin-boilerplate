package types

import (
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
)

type (
	CreateUserInput struct {
		UserID   uint
		FullName string
		Roles    []uint
		Email    string
		Status   enums.Status
		Password string
	}

	UpdateUserInput struct {
		ID       uint
		UserID   uint
		Status   enums.Status
		Password string
	}

	GetUserInput struct {
		ID     uint
		UserID uint
	}

	DeleteUserInput struct {
		ID     uint
		UserID uint
	}
)

type ListUsersInput struct {
	Search  string `json:"search"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}
