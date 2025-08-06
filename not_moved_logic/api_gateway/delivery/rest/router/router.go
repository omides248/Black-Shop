package router

import (
	pbcatalog "black-shop/api/proto/catalog/v1"
	pbidentity "black-shop/api/proto/identity/v1"
	"black-shop/internal/api_gateway/delivery/graphql/graph"
	"black-shop/internal/api_gateway/delivery/rest/handlers"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
)

func SetupRest(e *echo.Echo, catalogClient pbcatalog.CatalogServiceClient, identityClient pbidentity.IdentityServiceClient) {
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
}

func SetupGraphQL(e *echo.Echo, catalogClient pbcatalog.CatalogServiceClient) {

	schema := graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CatalogClient: catalogClient,
		},
	})

	srv := handler.New(schema)

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})

	srv.Use(extension.Introspection{})

	e.POST("/query", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/playground", func(c echo.Context) error {
		playground.Handler("GraphQL Playground", "/query").ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
