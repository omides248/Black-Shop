// file: backend/cmd/catalog_service/main.go
package main

import (
	"context"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "black-shop-service/api/proto/v1"
	"black-shop-service/internal/adapters/storage/mongodb"
	"black-shop-service/internal/app/catalog"
	"black-shop-service/pkg/config"
	"black-shop-service/pkg/logger"
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
		appLogger.Fatal("server failed to run", zap.Error(err))
	}
}

func run(ctx context.Context, cfg *config.Config, appLogger *zap.Logger) error {

	db, err := setupDatabase(ctx, cfg.MongoURI, appLogger)
	if err != nil {
		appLogger.Error("failed to setup database", zap.Error(err))
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(db.Client(), ctx)

	mongoRepo := mongodb.NewProductRepository(db, appLogger)
	catalogService := catalog.NewService(mongoRepo, appLogger)
	grpcServerHandler := catalog.NewGRPCServer(catalogService, appLogger)

	errCh := make(chan error, 1)

	go func() {
		errCh <- runGRPCServer(cfg.CatalogGRPCPort, grpcServerHandler, appLogger)
	}()

	go func() {
		errCh <- runRESTGateway(ctx, cfg.CatalogHTTPPort, cfg.CatalogGRPCAddr, appLogger)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("a server failed: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func setupDatabase(ctx context.Context, uri string, appLogger *zap.Logger) (*mongo.Database, error) {
	appLogger.Info("connecting to MongoDB...")

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

func runRESTGateway(ctx context.Context, httpPort, grpcAddr string, appLogger *zap.Logger) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterCatalogServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	handler := cors.AllowAll().Handler(mux)

	appLogger.Info("HTTP REST Gateway is running", zap.String("port", httpPort))
	return http.ListenAndServe(httpPort, handler)
}
