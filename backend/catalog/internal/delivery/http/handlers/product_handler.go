package handlers

import (
	"catalog/internal/application"
	"catalog/internal/delivery/http/dto"
	"catalog/internal/domain"
	"net/http"
	"pkg/echo/filter"
	"pkg/echo/pagination"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PaginatedResponse struct {
	Count    int64       `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  interface{} `json:"results"`
}

var productFilterSet = &filter.FilterSet{
	FilterFields: map[string]string{
		"category": "categoryId",
		"brand":    "brandId",
	},
	SearchFields:   []string{"name", "description"},
	OrderingFields: []string{"price", "created_at"},
}

type ProductHandler struct {
	service application.ProductService
	logger  *zap.Logger
}

func NewProductHandler(service application.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger.Named("product_http_handler"),
	}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req dto.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	productDomain, err := h.service.CreateProduct(c.Request().Context(), req.Name)
	if err != nil {
		h.logger.Error("failed to create product", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"product": toProductResponse(productDomain)})

}

func (h *ProductHandler) ListProducts(c echo.Context) error {

	page := c.Get("page").(int)
	limit := c.Get("size").(int)

	queryResult := productFilterSet.BuildMongoQuery(c)

	products, total, err := h.service.FindAllProducts(c.Request().Context(), queryResult.FilterQuery, queryResult.SortOptions, page, limit)
	if err != nil {
		h.logger.Error("failed to find all products", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "failed to retrieve products"})
	}

	results := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		results[i] = toProductResponse(p)
	}

	response := pagination.NewResponse(c, total, results)

	return c.JSON(http.StatusOK, response)
}

func toProductResponse(product *domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:   string(product.ID),
		Name: product.Name,
	}
}
