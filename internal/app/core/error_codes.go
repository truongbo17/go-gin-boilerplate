package core

// Stable application error codes shared by use cases and HTTP adapters.
const (
	ErrBadRequest          = 400
	ErrUnauthorized        = 401
	ErrForbidden           = 403
	ErrNotFound            = 404
	ErrTooManyRequests     = 429
	ErrUnprocessableEntity = 422
	ErrInternalServerError = 500

	ErrApp = 1000
)
