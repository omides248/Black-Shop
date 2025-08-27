package router

import (
	"catalog/config"
	"catalog/internal/application"
	"catalog/internal/delivery/http/handlers"
	"pkg/echo/pagination"
	"pkg/minio"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Setup(e *echo.Echo, categoryService application.CategoryService, productService application.ProductService, minioService *minio.Service, cfg *config.Config, logger *zap.Logger) {

	//e.Static(getStaticFilesPrefix(cfg.LocalStorage.StaticFilesPrefix), cfg.LocalStorage.PublicStoragePath)

	categoryHandler := handlers.NewCategoryHandler(categoryService, minioService, cfg, logger)
	productHandler := handlers.NewProductHandler(productService, logger)

	v1 := e.Group("/v1")
	{
		categories := v1.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.PATCH("/:id", categoryHandler.UpdateCategory)
			categories.GET("", categoryHandler.ListCategories)
		}

		products := v1.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.ListProducts, pagination.New())
		}

		e.GET("/images/:slug", categoryHandler.GetImage)
	}
}

//func getStaticFilesPrefix(staticFilesPrefix string) string {
//	if !strings.HasPrefix(staticFilesPrefix, "/") {
//		staticFilesPrefix = fmt.Sprintf("/%s", staticFilesPrefix)
//	}
//
//	return staticFilesPrefix
//}
