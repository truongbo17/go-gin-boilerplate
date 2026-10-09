package response

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/truongbo17/go-gin-boilerplate/internal/page"
)

// PageMeta adds HTTP links to a transport-independent page result.
func PageMeta[T any](ctx *gin.Context, result page.Result[T]) MetaData {
	if result.Number < 1 {
		result.Number = 1
	}
	if result.Size < 1 {
		result.Size = 20
	}
	path := strings.TrimRight(ctx.GetString("app_url"), "/") + ctx.Request.URL.Path
	lastPage := int((result.Total + int64(result.Size) - 1) / int64(result.Size))
	from := (result.Number-1)*result.Size + 1
	to := (result.Number-1)*result.Size + len(result.Items)
	if len(result.Items) == 0 {
		from = 0
		to = 0
	}
	meta := MetaData{
		CurrentPage:  result.Number,
		FirstPageUrl: pageURL(path, ctx.Request.URL.Query(), 1),
		From:         from,
		LastPage:     lastPage,
		LastPageUrl:  pageURL(path, ctx.Request.URL.Query(), lastPage),
		Path:         path,
		PerPage:      result.Size,
		To:           to,
		Total:        int(result.Total),
	}
	if result.Number < lastPage {
		meta.NextPageUrl = pageURL(path, ctx.Request.URL.Query(), result.Number+1)
	}
	if result.Number > 1 {
		meta.PrevPageUrl = pageURL(path, ctx.Request.URL.Query(), result.Number-1)
	}
	return meta
}

func pageURL(path string, query url.Values, number int) string {
	if number < 1 {
		return ""
	}
	parsed, err := url.Parse(path)
	if err != nil {
		return ""
	}
	copyQuery := url.Values{}
	for key, values := range query {
		copyQuery[key] = append([]string(nil), values...)
	}
	copyQuery.Set("page", strconv.Itoa(number))
	parsed.RawQuery = copyQuery.Encode()
	return parsed.String()
}
