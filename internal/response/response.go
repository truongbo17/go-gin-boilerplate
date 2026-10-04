package response

import (
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/i18n"
	"github.com/truongbo17/go-gin-boilerplate/internal/infra/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type BaseResponse struct {
	Status     bool        `json:"status"`
	StatusCode int         `json:"status_code"`
	RequestId  string      `json:"request_id"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	ExtraData  interface{} `json:"extra_data"`
	Error      interface{} `json:"errors"`
}

type PaginateResponse[T any] struct {
	Data *[]T `json:"data"`
	MetaData
}

type MetaData struct {
	CurrentPage  int    `json:"current_page"`
	FirstPageUrl string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageUrl  string `json:"last_page_url"`
	NextPageUrl  string `json:"next_page_url"`
	Path         string `json:"path"`
	PerPage      int    `json:"per_page"`
	PrevPageUrl  string `json:"prev_page_url"`
	To           int    `json:"to"`
	Total        int    `json:"total"`
}

func ReturnSuccess(ctx *gin.Context, data interface{}, extraData interface{}) {
	ctx.JSON(http.StatusOK, BaseResponse{
		Status:    true,
		RequestId: ctx.GetString(config.HeaderRequestID),
		Message:   i18n.GetMessage(ctx.GetString(config.HeaderLanguage), "success.200", nil),
		Data:      data,
		ExtraData: extraData,
	})
	ctx.Abort()
}

func ReturnError(ctx *gin.Context, status int, code int, err interface{}) {
	if err != nil && code >= ErrApp {
		logger.LogrusLogger.Errorf("Error: %+v, code: %d, path: %s, method: %s, request_id: %s",
			err,
			code,
			ctx.Request.URL.Path,
			ctx.Request.Method,
			ctx.GetString(config.HeaderRequestID),
		)
	}
	var errMsg interface{}
	if err != nil {
		switch v := err.(type) {
		case error:
			span := trace.SpanFromContext(ctx.Request.Context())
			span.RecordError(v)
			//errMsg = v.Error()
			errMsg = nil
		//case []byte:
		//	var jsonMap map[string]interface{}
		//	if err := json.Unmarshal(v, &jsonMap); err == nil {
		//		errMsg = jsonMap
		//	} else {
		//		errMsg = string(v)
		//	}
		case []string:
			if status == http.StatusUnprocessableEntity {
				errMsg = splitValidationErrors(v)
			} else {
				//errMsg = v
				errMsg = nil
			}
		case string:
			errMsg = v
		default:
			//errMsg = v
			errMsg = nil
		}
	}

	ctx.JSON(status, BaseResponse{
		RequestId:  ctx.GetString(config.HeaderRequestID),
		StatusCode: code,
		Message:    GetMessageError(ctx.GetString(config.HeaderLanguage), code),
		Error:      errMsg,
	})
	ctx.Abort()
}

func splitValidationErrors(errorsList []string) map[string]string {
	result := make(map[string]string, len(errorsList))

	for _, errString := range errorsList {
		errors := strings.Split(errString, ";")
		for _, e := range errors {
			e = strings.TrimSpace(e)
			if e != "" {
				parts := strings.SplitN(e, ":", 2)
				if len(parts) == 2 {
					field := strings.TrimSpace(parts[0])
					result[field] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return result
}

const (
	ErrBadRequest          = 400
	ErrUnauthorized        = 401
	ErrForbidden           = 403
	ErrNotFound            = 404
	ErrTooManyRequests     = 429
	ErrUnprocessableEntity = 422
	ErrInternalServerError = 500

	ErrApp = 1000

	ErrAuthLoginFailed       = 2000
	ErrAuthUserNotFound      = 2001
	ErrAuthWrongPassword     = 2002
	ErrAuthGenerateToken     = 2003
	ErrAuthUserExists        = 2004
	ErrAuthLogoutFailed      = 2005
	ErrAuthRegisterFailed    = 2006
	ErrUserListInternalError = 2007

	ErrChangePass   = 13008
	ErrRefreshToken = 13009

	ErrFindUserFailed   = 13010
	ErrUserNotFound     = 13011
	ErrUpdateUserFailed = 13012

	ErrRoleInternalError    = 14000
	ErrRoleNotFound         = 14001
	ErrRoleCreateFailed     = 14002
	ErrRoleUpdateFailed     = 14003
	ErrRoleDeleteFailed     = 14004
	ErrRoleAssignPermission = 14005
	ErrRoleAssignUser       = 14006

	ErrPermissionInternalError = 15000
	ErrPermissionNotFound      = 15001
	ErrPermissionCreateFailed  = 15002
	ErrPermissionUpdateFailed  = 15003
	ErrPermissionDeleteFailed  = 15004
)

var messageError = map[int]string{
	ErrBadRequest:          "error.400",
	ErrUnauthorized:        "error.401",
	ErrForbidden:           "error.403",
	ErrNotFound:            "error.404",
	ErrTooManyRequests:     "error.429",
	ErrUnprocessableEntity: "error.422",

	ErrInternalServerError: "error.500",

	ErrApp: "error.app",

	ErrAuthLoginFailed:    "auth.login_failed",
	ErrAuthRegisterFailed: "auth.register_failed",
	ErrAuthUserNotFound:   "auth.user_not_found",
	ErrAuthWrongPassword:  "auth.wrong_password",
	ErrAuthGenerateToken:  "auth.generate_token_failed",
	ErrAuthUserExists:     "auth.user_exists",
	ErrAuthLogoutFailed:   "auth.logout_failed",

	ErrUserListInternalError: "user_list.internal_error",

	ErrChangePass:   "err.user_change_pass",
	ErrRefreshToken: "err.refresh_token",

	ErrFindUserFailed:   "error.find_user_failed",
	ErrUserNotFound:     "error.user_not_found",
	ErrUpdateUserFailed: "error.update_user_failed",

	ErrRoleInternalError:    "role.internal_error",
	ErrRoleNotFound:         "role.not_found",
	ErrRoleCreateFailed:     "role.create_failed",
	ErrRoleUpdateFailed:     "role.update_failed",
	ErrRoleDeleteFailed:     "role.delete_failed",
	ErrRoleAssignPermission: "role.assign_permission_failed",
	ErrRoleAssignUser:       "role.assign_user_failed",

	ErrPermissionInternalError: "permission.internal_error",
	ErrPermissionNotFound:      "permission.not_found",
	ErrPermissionCreateFailed:  "permission.create_failed",
	ErrPermissionUpdateFailed:  "permission.update_failed",
	ErrPermissionDeleteFailed:  "permission.delete_failed",
}

func GetMessageError(lang string, statusCode int) string {
	msg := messageError[ErrApp]
	if v, ok := messageError[statusCode]; ok {
		msg = v
	}
	return i18n.GetMessage(lang, msg, nil)
}
