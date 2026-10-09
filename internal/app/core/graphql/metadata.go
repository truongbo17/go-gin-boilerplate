package graphql

import (
	"context"
	"errors"
	"strings"

	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
)

var (
	ErrUnsupportedKey = errors.New("unsupported entity key")
	ErrForbidden      = errors.New("entity access denied")
	ErrInvalidInput   = errors.New("invalid entity options input")
)

type EnumOption struct {
	Value int    `json:"value"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type EnumGroup struct {
	Key     string       `json:"key"`
	Options []EnumOption `json:"options"`
}

// Catalog returns copies so callers cannot change the shared enum definitions.
func Catalog() []EnumGroup {
	return []EnumGroup{{Key: "user_status", Options: []EnumOption{
		{int(enums.StatusInActive), "INACTIVE", "Inactive"},
		{int(enums.StatusActive), "ACTIVE", "Active"},
		{int(enums.StatusSuspended), "SUSPENDED", "Suspended"},
		{int(enums.StatusBanned), "BANNED", "Banned"},
	}}}
}

func EnumByKey(key string) (EnumGroup, bool) {
	for _, group := range Catalog() {
		if group.Key == strings.TrimSpace(key) {
			return group, true
		}
	}
	return EnumGroup{}, false
}

type EntityOption struct {
	ID    string  `json:"id"`
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Label string  `json:"label"`
	Extra *string `json:"extra"`
}

type PageMeta struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	LastPage int `json:"last_page"`
	Total    int `json:"total"`
}

type EntityPage struct {
	Key     string         `json:"key"`
	Options []EntityOption `json:"options"`
	Meta    PageMeta       `json:"meta"`
}

type EntityStore interface {
	ListOptions(context.Context, string, string, int, int) ([]EntityOption, int64, error)
}

type PermissionChecker interface {
	CheckPermission(context.Context, uint, string) (bool, error)
}

type Service struct {
	Store       EntityStore
	Permissions PermissionChecker
}

var entityPermissions = map[string]string{
	"users": "user:index", "roles": "role:index", "permissions": "permission:index",
}

func (service Service) EntityOptions(ctx context.Context, userID uint, key, keyword string, page, perPage int) (EntityPage, error) {
	key = strings.TrimSpace(key)
	permission, ok := entityPermissions[key]
	if !ok {
		return EntityPage{}, ErrUnsupportedKey
	}
	allowed, err := service.Permissions.CheckPermission(ctx, userID, permission)
	if err != nil {
		return EntityPage{}, err
	}
	if !allowed {
		return EntityPage{}, ErrForbidden
	}
	if page < 1 || perPage < 1 || perPage > 100 {
		return EntityPage{}, ErrInvalidInput
	}
	if len(keyword) > 100 {
		return EntityPage{}, ErrInvalidInput
	}
	options, total, err := service.Store.ListOptions(ctx, key, strings.TrimSpace(keyword), page, perPage)
	if err != nil {
		return EntityPage{}, err
	}
	return EntityPage{Key: key, Options: options, Meta: PageMeta{Page: page, PerPage: perPage, LastPage: max(1, int((total+int64(perPage)-1)/int64(perPage))), Total: int(total)}}, nil
}
