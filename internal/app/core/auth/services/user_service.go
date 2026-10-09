package services

import (
	"context"

	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	authcore "github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/enums"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	"github.com/truongbo17/go-gin-boilerplate/internal/page"
)

type UserService struct {
	UserRepository UserServiceRepository
}

type UserServiceRepository interface {
	CreateWithAdminRole(context.Context, *models.User, bool) error
	ListPage(context.Context, page.Query) (*page.Result[models.User], error)
	FindByID(context.Context, uint) (*models.User, error)
	Save(context.Context, *models.User) error
	GetUserRoles(context.Context, uint) ([]models.Role, error)
	AssignRoleToUser(context.Context, uint, uint, []uint) error
}

func NewUserService(userRepository UserServiceRepository) UserService {
	return UserService{UserRepository: userRepository}
}

func (us *UserService) CreateUser(ctx context.Context, input types.ProvisionUserInput) error {
	password, err := generatePassword(input.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(password),
		Status:   enums.StatusActive,
	}
	return us.UserRepository.CreateWithAdminRole(ctx, user, input.Admin)
}

func (us *UserService) ListUsers(ctx context.Context, input types.ListUsersInput) (*page.Result[models.User], *core.ErrorReturn) {
	users, err := us.UserRepository.ListPage(ctx, page.Query{Number: input.Page, Size: input.PerPage, Search: input.Search})
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrUserListInternalError,
			Err:       err,
		}
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, input types.UpdateUserInput) (*models.User, *core.ErrorReturn) {
	user, err := s.UserRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, &core.ErrorReturn{
			Err:       err,
			ErrorCode: authcore.ErrFindUserFailed,
		}
	}
	if user == nil {
		return nil, &core.ErrorReturn{
			Err:       nil,
			ErrorCode: authcore.ErrUserNotFound,
		}
	}
	pwd := user.Password
	if input.Password != "" {
		newPass, err := generatePassword(input.Password)
		if err != nil {
			return nil, &core.ErrorReturn{
				Err:       nil,
				ErrorCode: authcore.ErrUpdateUserFailed,
			}
		}
		pwd = string(newPass)
	}

	user.Status = input.Status
	user.Password = pwd

	if err = s.UserRepository.Save(ctx, user); err != nil {
		return nil, &core.ErrorReturn{
			Err:       err,
			ErrorCode: authcore.ErrUpdateUserFailed,
		}
	}

	return user, nil
}

func (us *UserService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, *core.ErrorReturn) {
	roles, err := us.UserRepository.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: authcore.ErrRoleInternalError,
			Err:       err,
		}
	}
	if roles == nil {
		roles = []models.Role{}
	}
	return roles, nil
}

func (us *UserService) AssignRoleToUser(ctx context.Context, input types.AssignRoleToUserInput) *core.ErrorReturn {
	if err := us.UserRepository.AssignRoleToUser(ctx, input.ActorID, input.UserID, input.RoleIDs); err != nil {
		return &core.ErrorReturn{ErrorCode: authcore.ErrRoleAssignUser, Err: err}
	}
	return nil
}
