package router

import (
	pbcatalog "black-shop/api/proto/catalog/v1"
	pbidentity "black-shop/api/proto/identity/v1"
	"black-shop/internal/api_gateway/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(catalogClient pbcatalog.CatalogServiceClient, identityClient pbidentity.IdentityServiceClient) *echo.Echo {
	e := echo.New()

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	catalogHandler := handlers.NewCatalogHandler(catalogClient)
	//identityHandler := handlers.NewIdentityHandler(identityClient)

	v1 := e.Group("/v1")
	{
		categories := v1.Group("/categories")
		{
			categories.POST("", catalogHandler.CreateCategory)
			categories.GET("", catalogHandler.ListCategories)
		}

		//auth := v1.Group("/auth")
		//{
		//	auth.POST("/register", identityHandler.Register)
		//}
	}
	return e
}
