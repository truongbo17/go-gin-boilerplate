package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core"
	"github.com/truongbo17/go-gin-boilerplate/internal/i18n"
	"go.opentelemetry.io/otel/trace"
)

type BaseResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"status_code"`
	RequestId  string `json:"request_id"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
	ExtraData  any    `json:"extra_data"`
	Error      any    `json:"errors"`
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

func ReturnSuccess(ctx *gin.Context, data any, extraData any) {
	ReturnSuccessStatus(ctx, http.StatusOK, i18n.GetMessage(ctx.GetString(config.HeaderLanguage), "success.200", nil), data, extraData)
}

func ReturnSuccessStatus(ctx *gin.Context, status int, message string, data any, extraData any) {
	ctx.JSON(status, BaseResponse{
		Status:     true,
		StatusCode: status,
		RequestId:  ctx.GetString(config.HeaderRequestID),
		Message:    message,
		Data:       data,
		ExtraData:  extraData,
	})
	ctx.Abort()
}

func ReturnError(ctx *gin.Context, status int, code int, message string, err any) {
	if err != nil && code >= core.ErrApp {
		if cause, ok := err.(error); ok {
			_ = ctx.Error(cause)
		}
	}
	var errMsg any
	if err != nil {
		switch v := err.(type) {
		case error:
			span := trace.SpanFromContext(ctx.Request.Context())
			span.RecordError(v)
		case []string:
			if status == http.StatusUnprocessableEntity {
				errMsg = splitValidationErrors(v)
			}
		case string:
			errMsg = v
		}
	}

	ctx.JSON(status, BaseResponse{
		RequestId:  ctx.GetString(config.HeaderRequestID),
		StatusCode: code,
		Message:    message,
		Error:      errMsg,
	})
	ctx.Abort()
}

func splitValidationErrors(errorsList []string) map[string]string {
	result := make(map[string]string, len(errorsList))

	for _, errString := range errorsList {
		errors := strings.SplitSeq(errString, ";")
		for e := range errors {
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
