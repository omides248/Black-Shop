package main

import (
	pbcatalog "black-shop/api/proto/catalog/v1"
	pbidentity "black-shop/api/proto/identity/v1"
	"black-shop/internal/api_gateway/delivery/rest/router"
	"black-shop/pkg/config"
	"black-shop/pkg/logger"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appLogger := logger.New(cfg.AppEnv)
	defer func(appLogger *zap.Logger) {
		_ = appLogger.Sync()
	}(appLogger)

	if err := run(context.Background(), cfg, appLogger); err != nil {
		appLogger.Fatal("api-gateway failed to run", zap.Error(err))
	}
}

func run(ctx context.Context, cfg *config.Config, appLogger *zap.Logger) error {

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Connect to catalog service
	catalogConn, err := grpc.NewClient(cfg.CatalogGRPCAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to catalog service: %w", err)
	}
	defer func(catalogConn *grpc.ClientConn) {
		_ = catalogConn.Close()
	}(catalogConn)
	catalogClient := pbcatalog.NewCatalogServiceClient(catalogConn)

	// Connect to identity service
	identityConn, err := grpc.NewClient(cfg.IdentityGRPCAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to identity service: %w", err)
	}
	defer func(identityConn *grpc.ClientConn) {
		_ = identityConn.Close()
	}(identityConn)
	identityClient := pbidentity.NewIdentityServiceClient(identityConn)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register REST routes
	router.SetupRest(e, catalogClient, identityClient)

	// Register GraphQL routes
	router.SetupGraphQL(e, catalogClient)

	httpPort := ":8080"
	appLogger.Info("Echo API Gateway is running on", zap.String("port", httpPort))
	return e.Start(httpPort)
}
