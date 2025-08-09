package handlers

import (
	"catalog/internal/application"
	"catalog/internal/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CreateCategoryRequest struct {
	Name     string  `json:"name"`
	ImageUrl *string `json:"image_url,omitempty"`
	ParentID *string `json:"parent_id,omitempty"`
}

type CategoryResponse struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	ImageURL      *string             `json:"imageUrl,omitempty"`
	ParentID      *string             `json:"parentId,omitempty"`
	Depth         int                 `json:"depth"`
	Subcategories []*CategoryResponse `json:"subcategory,omitempty"`
	CreatedAt     time.Time           `json:"createdAt"`
}

type CategoryHandler struct {
	service application.CategoryService
	logger  *zap.Logger
}

func NewCategoryHandler(service application.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		service: service,
		logger:  logger.Named("category_http_handler"),
	}
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	domainCategory, err := h.service.CreateCategory(c.Request().Context(), req.Name, req.ImageUrl, req.ParentID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"category": toCategoryResponse(domainCategory)})

}

func (h *CategoryHandler) ListCategories(c echo.Context) error {
	domainCategories, err := h.service.GetAllCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	treeCategories := buildCategoryTree(domainCategories)

	return c.JSON(http.StatusOK, echo.Map{"categories": treeCategories})
}

func buildCategoryTree(categories []*domain.Category) []*CategoryResponse {
	if len(categories) == 0 {
		return []*CategoryResponse{}
	}

	// Create Map
	categoryMap := make(map[domain.CategoryID]*CategoryResponse)
	for _, cat := range categories {
		dto := &CategoryResponse{
			ID:            string(cat.ID),
			Name:          cat.Name,
			ImageURL:      cat.ImageURL,
			Depth:         cat.Depth,
			CreatedAt:     cat.CreatedAt,
			Subcategories: []*CategoryResponse{},
		}
		if cat.ParentID != nil {
			pID := string(*cat.ParentID)
			dto.ParentID = &pID
		}
		categoryMap[cat.ID] = dto
	}

	var rootCategories []*CategoryResponse
	for _, cat := range categories {
		node := categoryMap[cat.ID]

		if cat.ParentID == nil {
			// Root node
			rootCategories = append(rootCategories, node)
		} else {
			// Child node
			parent, exists := categoryMap[*cat.ParentID] // Child's parent node

			if exists {
				parent.Subcategories = append(parent.Subcategories, node)
			}
		}
	}

	return rootCategories
}

func toCategoryResponse(cat *domain.Category) CategoryResponse {
	var parentID *string
	if cat.ParentID != nil {
		pID := string(*cat.ParentID)
		parentID = &pID
	}

	return CategoryResponse{
		ID:        string(cat.ID),
		Name:      cat.Name,
		ImageURL:  cat.ImageURL,
		ParentID:  parentID,
		Depth:     cat.Depth,
		CreatedAt: cat.CreatedAt,
	}
}
