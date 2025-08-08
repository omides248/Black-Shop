package router

import (
	"catalog/internal/application"
	"catalog/internal/delivery/http/handlers"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"pkg/pagination"
)

func Setup(e *echo.Echo, categoryService application.CategoryService, productService application.ProductService, logger *zap.Logger) {
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)
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
