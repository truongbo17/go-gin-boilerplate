package services

import (
	"context"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/repositories"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/types"
	baseModel "github.com/truongbo17/go-gin-boilerplate/internal/model"
	"github.com/truongbo17/go-gin-boilerplate/internal/repository"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
)

type UserService struct {
	UserRepository repositories.UserRepository
}

func NewUserService() UserService {
	return UserService{UserRepository: repositories.NewUserRepository()}
}

func (us *UserService) ListUsers(ctx context.Context, input types.ListUsersInput) (*response.PaginateResponse[models.User], *core.ErrorReturn) {
	params := repository.NewPaginateParam()
	params.Page = input.Page
	params.PerPage = input.PerPage
	params.Path = "/api/v1/rbac/users"

	if input.Search != "" {
		params.ExtraWheres = append(params.ExtraWheres, repository.WhereClause{
			Query: "username LIKE ? OR email LIKE ?",
			Args:  []interface{}{"%" + input.Search + "%", "%" + input.Search + "%"},
		})
	}

	users, err := us.UserRepository.Paginate(ctx, params)
	if err != nil {
		return nil, &core.ErrorReturn{
			ErrorCode: response.ErrUserListInternalError,
			Err:       err,
		}
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, input types.UpdateUserInput) (*models.User, *core.ErrorReturn) {
	user, err := s.UserRepository.FindOneByCondition(ctx, models.User{
		BasicModel: baseModel.BasicModel{
			BasicIDModel: baseModel.BasicIDModel{
				ID: input.ID,
			},
		},
	})
	if err != nil {
		return nil, &core.ErrorReturn{
			Err:       err,
			ErrorCode: response.ErrFindUserFailed,
		}
	}
	if user == nil {
		return nil, &core.ErrorReturn{
			Err:       nil,
			ErrorCode: response.ErrUserNotFound,
		}
	}
	serviceAuth := NewAuthService()

	pwd := user.Password
	if input.Password != "" {
		newPass, err := serviceAuth.GeneratePassword(input.Password)
		if err != nil {
			return nil, &core.ErrorReturn{
				Err:       nil,
				ErrorCode: response.ErrUpdateUserFailed,
			}
		}
		pwd = string(newPass)
	}

	user.Status = input.Status
	user.Password = pwd

	if err = s.UserRepository.Save(ctx, user); err != nil {
		return nil, &core.ErrorReturn{
			Err:       err,
			ErrorCode: response.ErrUpdateUserFailed,
		}
	}

	return user, nil
}
