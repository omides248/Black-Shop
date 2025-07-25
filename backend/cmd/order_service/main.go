package main

import (
	"context"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "black-shop-service/api/proto/v1"
	"black-shop-service/internal/adapters/storage/postgresql"
	redisStorage "black-shop-service/internal/adapters/storage/redis"
	"black-shop-service/internal/app/order"
	"black-shop-service/pkg/auth"
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
		appLogger.Fatal("order-service failed to run", zap.Error(err))
	}
}

func run(ctx context.Context, cfg *config.Config, appLogger *zap.Logger) error {
	// --- Database Connections ---
	appLogger.Info("connecting to Postgresql...")
	pgConn, err := pgx.Connect(ctx, cfg.PostgresOrderURI)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer func(pgConn *pgx.Conn, ctx context.Context) {
		_ = pgConn.Close(ctx)
	}(pgConn, ctx)
	appLogger.Info("Successfully connected to PostgreSQL")

	appLogger.Info("connecting to Redis...")
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	defer func(redisClient *redis.Client) {
		_ = redisClient.Close()
	}(redisClient)
	appLogger.Info("Successfully connected to Redis")

	// --- Layers ---
	tokenManager := auth.NewTokenManager("a_very_secret_key")

	cartRepo := redisStorage.NewCartRepository(redisClient, appLogger)
	orderRepo, err := postgresql.NewOrderRepository(pgConn, appLogger)
	if err != nil {
		return fmt.Errorf("failed to create order repository: %w", err)
	}

	orderService := order.NewService(cartRepo, orderRepo, appLogger)
	grpcHandler := order.NewGRPCServer(orderService, appLogger)

	// --- Servers ---
	errCh := make(chan error, 1)
	go func() {
		errCh <- runGRPCServer(cfg.OrderGRPCPort, grpcHandler, tokenManager, appLogger)
	}()
	go func() {
		errCh <- runRESTGateway(ctx, cfg.OrderHTTPPort, cfg.OrderGRPCAddr, appLogger)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("a server failed: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func runGRPCServer(port string, handler pb.OrderServiceServer, tm *auth.TokenManager, appLogger *zap.Logger) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", port, err)
	}

	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(tm.AuthenticationInterceptor))
	pb.RegisterOrderServiceServer(gRPCServer, handler)

	appLogger.Info("Order gRPC Server is running", zap.String("port", port))
	return gRPCServer.Serve(lis)
}

func runRESTGateway(ctx context.Context, httpPort, grpcAddr string, appLogger *zap.Logger) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return fmt.Errorf("failed to register order gateway: %w", err)
	}

	handler := cors.AllowAll().Handler(mux)

	appLogger.Info("Order HTTP REST Gateway is running", zap.String("port", httpPort))
	return http.ListenAndServe(httpPort, handler)
}
