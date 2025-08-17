package router

import (
	"catalog/config"
	"catalog/internal/application"
	"catalog/internal/delivery/http/handlers"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"pkg/echo/pagination"
	"pkg/local_storage"
	"strings"
)

func Setup(e *echo.Echo, categoryService application.CategoryService, productService application.ProductService, localStorageService *local_storage.Service, cfg *config.Config, logger *zap.Logger) {

	e.Static(getStaticFilesPrefix(cfg.LocalStorage.StaticFilesPrefix), cfg.LocalStorage.PublicStoragePath)

	categoryHandler := handlers.NewCategoryHandler(categoryService, localStorageService, cfg, logger)
	productHandler := handlers.NewProductHandler(productService, logger)

	v1 := e.Group("/v1")
	{
		categories := v1.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.ListCategories)
		}

		products := v1.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.ListProducts, pagination.New())
		}

	}
}

func getStaticFilesPrefix(staticFilesPrefix string) string {
	if !strings.HasPrefix(staticFilesPrefix, "/") {
		staticFilesPrefix = fmt.Sprintf("/%s", staticFilesPrefix)
	}

	return staticFilesPrefix
}
