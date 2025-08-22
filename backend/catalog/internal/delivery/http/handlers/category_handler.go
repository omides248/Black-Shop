package handlers

import (
	"catalog/config"
	"catalog/internal/application"
	"catalog/internal/delivery/http/dto"
	"catalog/internal/domain"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"pkg/local_storage"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	service             application.CategoryService
	localStorageService *local_storage.Service
	config              *config.Config
	logger              *zap.Logger
}

func NewCategoryHandler(service application.CategoryService, localStorageService *local_storage.Service, cfg *config.Config, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		service:             service,
		localStorageService: localStorageService,
		config:              cfg,
		logger:              logger.Named("category_http_handler"),
	}
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var req dto.CreateCategoryRequest

	req.Name = c.FormValue("name")

	parentID := c.FormValue("parentId")
	if parentID != "" {
		req.ParentID = &parentID
	}

	if file, err := c.FormFile("image"); err == nil {
		req.Image = file
	}

	if err := req.Validate(); err != nil {
		return err
	}

	domainCategory, err := h.service.CreateCategory(c.Request().Context(), req.Name, nil, req.ParentID)
	if err != nil {
		return err
	}

	if req.Image != nil {
		src, err := req.Image.Open()
		if err != nil {
			h.logger.Error("failed to open uploaded file", zap.Error(err))
			return err
		}
		defer func(src multipart.File) {
			_ = src.Close()
		}(src)

		subDirectory := "categories"
		relativePath, err := h.localStorageService.UploadFile(subDirectory, req.Image.Filename, src)
		if err != nil {
			h.logger.Error("failed to upload file to MinIO", zap.Error(err))
			return err
		}

		domainCategory.Image = &relativePath
		err = h.service.UpdateCategory(c.Request().Context(), domainCategory)
		if err != nil {
			h.logger.Error("failed to update category with image URL", zap.Error(err))
			return err
		}
	} else if !errors.Is(err, http.ErrMissingBoundary) && !errors.Is(err, http.ErrNotMultipart) {
		h.logger.Warn("image file not provided", zap.Error(err))
	}

	response := h.toCategoryResponse(c, domainCategory)
	return c.JSON(http.StatusOK, echo.Map{"category": response})
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	categoryID := c.Param("id")

	var req dto.UpdateCategoryRequest

	name := c.FormValue("name")
	if name != "" {
		req.Name = &name
	}

	parentID := c.FormValue("parentId")
	if parentID != "" {
		req.ParentID = &parentID
	}

	if file, err := c.FormFile("image"); err == nil {
		req.Image = file
	}

	if err := req.Validate(); err != nil {
		return err
	}

	category, err := h.service.FindByID(c.Request().Context(), domain.CategoryID(categoryID))
	if err != nil {
		return err
	}
	oldImage := category.Image // Store old image path

	if req.Name != nil {
		category.Name = *req.Name
	}

	if req.ParentID != nil {
		parentCatID := domain.CategoryID(*req.ParentID)
		category.ParentID = &parentCatID
	}

	if req.Image != nil {
		src, err := req.Image.Open()
		if err != nil {
			h.logger.Error("failed to open uploaded file", zap.Error(err))
			return err
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {
				h.logger.Error("failed to close uploaded file", zap.Error(err))
			}
		}(src)

		subDirectory := "categories"
		relativePath, err := h.localStorageService.UploadFile(subDirectory, req.Image.Filename, src)
		if err != nil {
			h.logger.Error("failed to upload new image file", zap.Error(err))
			return err
		}
		category.Image = &relativePath
	}

	fmt.Printf("category %+v\n", category)
	fmt.Printf("category %+v\n", category)
	fmt.Printf("category %+v\n", category)
	fmt.Printf("category %+v\n", category)

	if err := h.service.UpdateCategory(c.Request().Context(), category); err != nil {
		return err
	}

	if req.Image != nil && oldImage != nil {
		err := h.localStorageService.DeleteFile(*oldImage)
		if err != nil {
			h.logger.Error("failed to delete old image file", zap.Error(err))
		}
	}

	response := h.toCategoryResponse(c, category)
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

func (h *CategoryHandler) buildCategoryTree(c echo.Context, categories []*domain.Category) []*dto.CategoryResponse {
	if len(categories) == 0 {
		return []*dto.CategoryResponse{}
	}

	// Create Map
	categoryMap := make(map[domain.CategoryID]*dto.CategoryResponse)
	for _, cat := range categories {
		categoryDto := h.toCategoryResponse(c, cat)

		if cat.ParentID != nil {
			pID := string(*cat.ParentID)
			categoryDto.ParentID = &pID
		}

		categoryMap[cat.ID] = &categoryDto
	}

	var rootCategories []*dto.CategoryResponse
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

func (h *CategoryHandler) toCategoryResponse(c echo.Context, cat *domain.Category) dto.CategoryResponse {
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

	return dto.CategoryResponse{
		ID:        string(cat.ID),
		Name:      cat.Name,
		Image:     image,
		ParentID:  parentID,
		Depth:     cat.Depth,
		CreatedAt: cat.CreatedAt,
	}
}
