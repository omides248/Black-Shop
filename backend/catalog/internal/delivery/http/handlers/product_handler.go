package handlers

import (
	"catalog/internal/application"
	"catalog/internal/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"pkg/pagination"
)

type CreateProductRequest struct {
	Name string `json:"name"  validate:"required,min=3"`
}

type ProductResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductHandler struct {
	service application.ProductService
	logger  *zap.Logger
}

type PaginatedResponse struct {
	Count    int64       `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  interface{} `json:"results"`
}

func NewProductHandler(service application.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger.Named("product_http_handler"),
	}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// TODO validate with thirty party library

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

	products, total, err := h.service.FindAllProducts(c.Request().Context(), page, limit)
	if err != nil {
		h.logger.Error("failed to find all products", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "failed to retrieve products"})
	}

	results := make([]ProductResponse, len(products))
	for i, p := range products {
		results[i] = toProductResponse(p)
	}

	response := pagination.NewResponse(c, total, results)

	return c.JSON(http.StatusOK, response)
}

func toProductResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:   string(product.ID),
		Name: product.Name,
	}
}
