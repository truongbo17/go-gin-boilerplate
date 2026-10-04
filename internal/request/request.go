package request

// BaseRequestPaginate holds the common query parameters for list endpoints.
type BaseRequestPaginate struct {
	Page    int `form:"page" json:"page"`
	PerPage int `form:"per_page" json:"per_page"`
}
