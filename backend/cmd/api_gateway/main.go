package main

import (
	pbcatalog "black-shop/api/proto/catalog/v1"
	pbidentity "black-shop/api/proto/identity/v1"
	"black-shop/internal/api_gateway/router"
	"black-shop/pkg/config"
	"black-shop/pkg/logger"
	"context"
	"fmt"
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
		appLogger.Fatal("catalog-service failed to run", zap.Error(err))
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

	r := router.Setup(catalogClient, identityClient)

	httpPort := ":8080"
	appLogger.Info("Echo API Gateway is running", zap.String("port", httpPort))
	return r.Start(httpPort)

}
