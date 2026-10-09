package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	authcore "github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth"
	"github.com/truongbo17/go-gin-boilerplate/internal/response"
)

type errorPresentation struct {
	status int
	key    string
}

var errorPresentations = map[int]errorPresentation{
	core.ErrBadRequest:          {http.StatusBadRequest, "error.400"},
	core.ErrUnauthorized:        {http.StatusUnauthorized, "error.401"},
	core.ErrForbidden:           {http.StatusForbidden, "error.403"},
	core.ErrNotFound:            {http.StatusNotFound, "error.404"},
	core.ErrTooManyRequests:     {http.StatusTooManyRequests, "error.429"},
	core.ErrUnprocessableEntity: {http.StatusUnprocessableEntity, "error.422"},
	core.ErrInternalServerError: {http.StatusInternalServerError, "error.500"},
	core.ErrApp:                 {http.StatusBadRequest, "error.app"},

	authcore.ErrAuthLoginFailed:         {http.StatusUnauthorized, "auth.login_failed"},
	authcore.ErrAuthUserNotFound:        {http.StatusUnauthorized, "auth.user_not_found"},
	authcore.ErrAuthWrongPassword:       {http.StatusUnauthorized, "auth.wrong_password"},
	authcore.ErrAuthGenerateToken:       {http.StatusInternalServerError, "auth.generate_token_failed"},
	authcore.ErrAuthUserExists:          {http.StatusConflict, "auth.user_exists"},
	authcore.ErrAuthLogoutFailed:        {http.StatusBadRequest, "auth.logout_failed"},
	authcore.ErrAuthRegisterFailed:      {http.StatusInternalServerError, "auth.register_failed"},
	authcore.ErrAuthResetUnavailable:    {http.StatusServiceUnavailable, "auth.reset_unavailable"},
	authcore.ErrAuthResetInvalid:        {http.StatusBadRequest, "auth.reset_invalid"},
	authcore.ErrAuthResetFailed:         {http.StatusInternalServerError, "auth.reset_failed"},
	authcore.ErrUserListInternalError:   {http.StatusInternalServerError, "user_list.internal_error"},
	authcore.ErrChangePass:              {http.StatusBadRequest, "err.user_change_pass"},
	authcore.ErrRefreshToken:            {http.StatusBadRequest, "err.refresh_token"},
	authcore.ErrFindUserFailed:          {http.StatusBadRequest, "error.find_user_failed"},
	authcore.ErrUserNotFound:            {http.StatusNotFound, "error.user_not_found"},
	authcore.ErrUpdateUserFailed:        {http.StatusBadRequest, "error.update_user_failed"},
	authcore.ErrRoleInternalError:       {http.StatusInternalServerError, "role.internal_error"},
	authcore.ErrRoleNotFound:            {http.StatusNotFound, "role.not_found"},
	authcore.ErrRoleCreateFailed:        {http.StatusBadRequest, "role.create_failed"},
	authcore.ErrRoleUpdateFailed:        {http.StatusBadRequest, "role.update_failed"},
	authcore.ErrRoleDeleteFailed:        {http.StatusBadRequest, "role.delete_failed"},
	authcore.ErrRoleAssignPermission:    {http.StatusBadRequest, "role.assign_permission_failed"},
	authcore.ErrRoleAssignUser:          {http.StatusBadRequest, "role.assign_user_failed"},
	authcore.ErrPermissionInternalError: {http.StatusInternalServerError, "permission.internal_error"},
	authcore.ErrPermissionNotFound:      {http.StatusNotFound, "permission.not_found"},
	authcore.ErrPermissionCreateFailed:  {http.StatusBadRequest, "permission.create_failed"},
	authcore.ErrPermissionUpdateFailed:  {http.StatusBadRequest, "permission.update_failed"},
	authcore.ErrPermissionDeleteFailed:  {http.StatusBadRequest, "permission.delete_failed"},
}

func ReturnError(ctx *gin.Context, failure *core.ErrorReturn) {
	presentation, ok := errorPresentations[failure.ErrorCode]
	if !ok {
		presentation = errorPresentations[core.ErrApp]
	}
	response.ReturnError(ctx, presentation.status, failure.ErrorCode, authMessage(ctx.GetString(config.HeaderLanguage), presentation.key), failure.Err)
}
