package services

import (
	"testing"

	"github.com/truongbo17/go-gin-boilerplate/internal/infra/database"
	"gorm.io/gorm"
)

func TestServiceConstructorsUseCurrentDatabase(t *testing.T) {
	previous := database.DB
	t.Cleanup(func() { database.DB = previous })

	first := &gorm.DB{}
	second := &gorm.DB{}
	database.DB = first
	_ = NewAuthService()
	_ = NewRoleService()
	_ = NewPermissionService()
	_ = NewUserService()

	database.DB = second
	if got := NewAuthService().UserRepository.DB; got != second {
		t.Fatalf("auth service kept stale database %p; want %p", got, second)
	}
	if got := NewRoleService().RoleRepository.DB; got != second {
		t.Fatalf("role service kept stale database %p; want %p", got, second)
	}
	if got := NewPermissionService().PermissionRepository.DB; got != second {
		t.Fatalf("permission service kept stale database %p; want %p", got, second)
	}
	if got := NewUserService().UserRepository.DB; got != second {
		t.Fatalf("user service kept stale database %p; want %p", got, second)
	}
}
