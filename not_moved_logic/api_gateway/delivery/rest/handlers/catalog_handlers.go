package handlers

import (
	"net/http"

	pbcatalog "black-shop/api/proto/catalog/v1"
	"github.com/labstack/echo/v4"
)

type CatalogHandler struct {
	client pbcatalog.CatalogServiceClient
}

func NewCatalogHandler(client pbcatalog.CatalogServiceClient) *CatalogHandler {
	return &CatalogHandler{client: client}
}

func (h *CatalogHandler) CreateCategory(c echo.Context) error {
	var req pbcatalog.CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	res, err := h.client.CreateCategory(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *CatalogHandler) ListCategories(c echo.Context) error {
	res, err := h.client.ListCategories(c.Request().Context(), &pbcatalog.ListCategoriesRequest{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
