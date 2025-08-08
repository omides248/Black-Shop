package paginator

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Count    int64       `json:"count"`
	Next     *string     `json:"next,omitempty"`
	Previous *string     `json:"previous,omitempty"`
	Results  interface{} `json:"results"`
}

func New(c echo.Context, totalItems int64, results interface{}) Response {
	page := c.Get("page").(int)
	limit := c.Get("size").(int)

	response := Response{
		Count:   totalItems,
		Results: results,
	}

	if (int64(page) * int64(limit)) < totalItems {
		nextURL := buildURL(c, page+1, limit)
		response.Next = &nextURL
	}

	if page > 1 {
		prevURL := buildURL(c, page-1, limit)
		response.Previous = &prevURL
	}

	return response
}

func buildURL(c echo.Context, page, limit int) string {
	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}

	queryParams := c.QueryParams()
	queryParams.Set("page", fmt.Sprintf("%d", page))
	queryParams.Set("limit", fmt.Sprintf("%d", limit))

	baseURL := fmt.Sprintf("%s://%s%s", scheme, c.Request().Host, c.Request().URL.Path)
	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}

//
//import (
//	"fmt"
//	"github.com/labstack/echo/v4"
//	"math"
//	"strconv"
//)
//
//const (
//	DefaultPageSize = 20
//	MaxPageSize     = 2
//	DefaultPage     = 1
//)
//
//type PaginatedResponse struct {
//	Count      int64       `json:"count"`
//	TotalPages int         `json:"total_pages"`
//	Next       *string     `json:"next"`
//	Previous   *string     `json:"previous"`
//	Results    interface{} `json:"results"`
//}
//
//func NewPaginator(c echo.Context, totalItems int64, pageSize int, results interface{}) PaginatedResponse {
//	page := ParseIntQuery(c, "page", DefaultPage)
//
//	totalPages := 0
//	if totalItems > 0 {
//		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
//	}
//
//	if page > totalPages && totalPages > 0 {
//		results = []interface{}{}
//	}
//
//	var nextURL, prevURL *string
//
//	if page < totalPages {
//		next := buildURL(c, page+1, pageSize)
//		nextURL = &next
//	}
//
//	if page > 1 && page <= totalPages {
//		prev := buildURL(c, page-1, pageSize)
//		prevURL = &prev
//	}
//
//	return PaginatedResponse{
//		Count:      totalItems,
//		TotalPages: totalPages,
//		Next:       nextURL,
//		Previous:   prevURL,
//		Results:    results,
//	}
//}
//
//func ParseIntQuery(c echo.Context, key string, defaultValue int) int {
//	valueStr := c.QueryParam(key)
//	if valueStr == "" {
//		return defaultValue
//	}
//
//	value, err := strconv.ParseInt(valueStr, 10, 64)
//	if err != nil || value <= 0 {
//		return defaultValue
//	}
//
//	if value > math.MaxInt32 {
//		return defaultValue
//	}
//	return int(value)
//}
//
//func buildURL(c echo.Context, page, pageSize int) string {
//	scheme := "http"
//	if c.Request().TLS != nil {
//		scheme = "https"
//	}
//
//	baseURL := fmt.Sprintf("%s://%s%s", scheme, c.Request().Host, c.Request().URL.Path)
//	return fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page, pageSize)
//}
