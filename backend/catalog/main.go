package main

import (
	"catalog/internal/adapters/storage/mongodb"
	"catalog/internal/application"
	"context"
	"fmt"
	"log"
	"net"
	"pkg/config"
	"pkg/logger"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "catalog/api/proto/v1"

	grpcserver "catalog/internal/delivery/grpc"
)

func main() {
	cfg, err := config.LoadConfig("./config")
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
	// Setup database
	db, err := setupDatabase(ctx, cfg.MongoURI, appLogger)
	if err != nil {
		appLogger.Error("failed to setup database", zap.Error(err))
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(db.Client(), ctx)

	// --- Repositories ---
	productRepo, err := mongodb.NewProductRepository(db, appLogger)
	if err != nil {
		return fmt.Errorf("failed to create product repository: %w", err)
	}

	categoryRepo, err := mongodb.NewCategoryRepository(db, appLogger)
	if err != nil {
		return fmt.Errorf("failed to create category repository: %w", err)
	}

	// --- Application Services ---
	productService := application.NewProductService(productRepo, appLogger)
	categoryService := application.NewCategoryService(categoryRepo, productRepo, appLogger)

	grpcServerDeps := grpcserver.ServerDependencies{
		ProductService:  productService,
		CategoryService: categoryService,
		Logger:          appLogger,
	}

	// --- Delivery Layer (gRPC Handler) ---
	grpcServerHandler := grpcserver.NewServer(grpcServerDeps, appLogger)

	// Setup Servers
	errCh := make(chan error, 1)
	go func() {
		errCh <- runGRPCServer(cfg.CatalogGRPCPort, grpcServerHandler, appLogger)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("a server failed: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func setupDatabase(ctx context.Context, uri string, appLogger *zap.Logger) (*mongo.Database, error) {
	appLogger.Info("connecting to MongoDB...", zap.String("uri", uri))
	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(uri)
	mongoClient, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}
	if err := mongoClient.Ping(connectCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	appLogger.Info("successfully connected to MongoDB")
	return mongoClient.Database("black_shop_db"), nil
}

func runGRPCServer(port string, handler pb.CatalogServiceServer, appLogger *zap.Logger) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", port, err)
	}
	gRPCServer := grpc.NewServer()
	pb.RegisterCatalogServiceServer(gRPCServer, handler)
	appLogger.Info("gRPC Server is running", zap.String("port", port))
	return gRPCServer.Serve(lis)
}
