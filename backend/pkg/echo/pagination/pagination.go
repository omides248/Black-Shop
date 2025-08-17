// Package pagination provides a middleware for the echo web framework to handle
// pagination. It allows for the usage of url parameters like `?page=1&size=25`
// to paginate data on your API. The values will be propagated throughout the
// request context.
package pagination

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetPage returns the page number from the request context.
func GetPage(c echo.Context, customOptions ...CustomOption) (int, error) {
	opts := applyCustomOptionsToDefault(customOptions...)
	page, ok := c.Get(opts.PageText).(int)
	if !ok {
		return 0, fmt.Errorf("%s not found in context, please ensure pagination middleware is used with the correct options", opts.PageText)
	}
	return page, nil
}

// GetPageSize returns the page size from the request context.
func GetPageSize(c echo.Context, customOptions ...CustomOption) (int, error) {
	opts := applyCustomOptionsToDefault(customOptions...)
	size, ok := c.Get(opts.SizeText).(int)
	if !ok {
		return 0, fmt.Errorf("%s not found in context, please ensure pagination middleware is used with the correct options", opts.SizeText)
	}
	return size, nil
}

// New returns a new pagination middleware with custom values.
func New(customOptions ...CustomOption) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			p := &paginator{
				opts: applyCustomOptionsToDefault(customOptions...),
				c:    c,
			}

			page, err := p.getPageFromRequest()
			if err != nil {
				return p.abortWithBadRequest(err)
			}

			if err := p.validatePage(page); err != nil {
				return p.abortWithBadRequest(err)
			}

			pageSize, err := p.getPageSizeFromRequest()
			if err != nil {
				return p.abortWithBadRequest(err)
			}

			if err := p.validatePageSize(pageSize); err != nil {
				return p.abortWithBadRequest(err)
			}

			p.setPageAndPageSize(page, pageSize)

			return next(c)
		}
	}
}

type paginator struct {
	opts options
	c    echo.Context
}

func (p *paginator) abortWithBadRequest(err error) error {
	return p.c.String(http.StatusBadRequest, err.Error())
}

func (p *paginator) getPageFromRequest() (int, error) {
	return p.getIntFromContextWithDefault(p.opts.PageText, p.opts.DefaultPage)
}

func (p *paginator) getPageSizeFromRequest() (int, error) {
	return p.getIntFromContextWithDefault(p.opts.SizeText, p.opts.DefaultPageSize)
}

func (p *paginator) getIntFromContextWithDefault(key string, defaultValue int) (int, error) {
	valueStr := p.c.QueryParam(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("%s parameter must be an integer", key)
	}

	return valueInt, nil
}

func (p *paginator) validatePage(page int) error {
	if page <= 0 {
		return fmt.Errorf("%s must be positive", p.opts.PageText)
	}

	return nil
}

func (p *paginator) validatePageSize(size int) error {
	if size < p.opts.MinPageSize || size > p.opts.MaxPageSize {
		return fmt.Errorf(
			"%s must be between %d and %d",
			p.opts.SizeText,
			p.opts.MinPageSize,
			p.opts.MaxPageSize,
		)
	}

	return nil
}

func (p *paginator) constructHeaderKey(key string) string {
	return p.opts.HeaderPrefix + key
}

func (p *paginator) setPageAndPageSize(page int, size int) {
	p.c.Set(p.opts.PageText, page)
	p.c.Set(p.opts.SizeText, size)

	p.c.Response().Header().Set(p.constructHeaderKey(p.opts.PageText), strconv.Itoa(page))
	p.c.Response().Header().Set(p.constructHeaderKey(p.opts.SizeText), strconv.Itoa(size))
}

type Response struct {
	Count      int64       `json:"count"`
	TotalPages int         `json:"total_pages"`
	Next       *string     `json:"next"`
	Previous   *string     `json:"previous"`
	Results    interface{} `json:"results"`
}

func NewResponse(c echo.Context, totalItems int64, results interface{}) Response {
	page, err := GetPage(c)
	if err != nil {
		page = 1
	}
	limit, err := GetPageSize(c)
	if err != nil {
		limit = 10
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(limit)))
	}

	response := Response{
		Count:      totalItems,
		TotalPages: totalPages,
		Results:    results,
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
	queryParams.Set("size", fmt.Sprintf("%d", limit))
	queryParams.Del("limit")

	baseURL := fmt.Sprintf("%s://%s%s", scheme, c.Request().Host, c.Request().URL.Path)
	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}
