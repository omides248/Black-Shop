package handlers

import (
	"catalog/config"
	"catalog/internal/application"
	"catalog/internal/domain"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"pkg/local_storage"
	"time"
)

type CategoryHandler struct {
	service             application.CategoryService
	localStorageService *local_storage.Service
	config              *config.Config
	logger              *zap.Logger
}

//type CreateCategoryRequest struct {
//	Name     string  `json:"name"`
//	ParentID *string `json:"parent_id,omitempty"`
//}

type CategoryResponse struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Image         *string             `json:"image,omitempty"`
	ParentID      *string             `json:"parentId,omitempty"`
	Depth         int                 `json:"depth"`
	Subcategories []*CategoryResponse `json:"subcategory,omitempty"`
	CreatedAt     time.Time           `json:"createdAt"`
}

func NewCategoryHandler(service application.CategoryService, localStorageService *local_storage.Service, cfg *config.Config, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		service:             service,
		localStorageService: localStorageService,
		config:              cfg,
		logger:              logger.Named("category_http_handler"),
	}
}

//func (h *CategoryHandler) CreateCategory(c echo.Context) error {
//	var req CreateCategoryRequest
//	if err := c.Bind(&req); err != nil {
//		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
//	}
//
//	domainCategory, err := h.service.CreateCategory(c.Request().Context(), req.Name, req.ParentID)
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
//	}
//
//	return c.JSON(http.StatusOK, echo.Map{"category": toCategoryResponse(domainCategory)})
//}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	name := c.FormValue("name")
	parentIDStr := c.FormValue("parent_id")

	var parentID *string
	if parentIDStr != "" {
		parentID = &parentIDStr
	}

	domainCategory, err := h.service.CreateCategory(c.Request().Context(), name, nil, parentID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	file, err := c.FormFile("image")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			h.logger.Error("failed to open uploaded file", zap.Error(err))
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to open uploaded file"})
		}
		defer func(src multipart.File) {
			_ = src.Close()
		}(src)

		subDirectory := "categories"
		relativePath, err := h.localStorageService.UploadFile(subDirectory, file.Filename, src)
		if err != nil {
			h.logger.Error("failed to upload file to MinIO", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to upload file"})
		}

		domainCategory.Image = &relativePath
		err = h.service.UpdateCategory(c.Request().Context(), domainCategory)
		if err != nil {
			h.logger.Error("failed to update category with image URL", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update category with image URL"})
		}
	} else if !errors.Is(err, http.ErrMissingBoundary) && !errors.Is(err, http.ErrNotMultipart) {
		h.logger.Warn("image file not provided", zap.Error(err))
	}

	response := h.toCategoryResponse(c, domainCategory)
	return c.JSON(http.StatusOK, echo.Map{"category": response})
}

func (h *CategoryHandler) ListCategories(c echo.Context) error {
	domainCategories, err := h.service.GetAllCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	treeCategories := h.buildCategoryTree(c, domainCategories)

	return c.JSON(http.StatusOK, echo.Map{"categories": treeCategories})
}

func (h *CategoryHandler) buildCategoryTree(c echo.Context, categories []*domain.Category) []*CategoryResponse {
	if len(categories) == 0 {
		return []*CategoryResponse{}
	}

	// Create Map
	categoryMap := make(map[domain.CategoryID]*CategoryResponse)
	for _, cat := range categories {
		dto := h.toCategoryResponse(c, cat)

		if cat.ParentID != nil {
			pID := string(*cat.ParentID)
			dto.ParentID = &pID
		}
		categoryMap[cat.ID] = &dto
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

func (h *CategoryHandler) toCategoryResponse(c echo.Context, cat *domain.Category) CategoryResponse {
	var parentID *string
	if cat.ParentID != nil {
		pID := string(*cat.ParentID)
		parentID = &pID
	}

	var image *string
	if cat.Image != nil {

		scheme := "https"
		if c.Request().TLS == nil {
			scheme = "http"
		}

		fullURL := fmt.Sprintf("%s://%s/%s/%s", scheme, h.config.General.Host, h.config.LocalStorage.StaticFilesPrefix, *cat.Image)
		image = &fullURL
	}

	return CategoryResponse{
		ID:        string(cat.ID),
		Name:      cat.Name,
		Image:     image,
		ParentID:  parentID,
		Depth:     cat.Depth,
		CreatedAt: cat.CreatedAt,
	}
}
