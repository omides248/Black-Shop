package main

import (
	"black-shop/pkg/auth"
	"context"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "black-shop/api/proto/v1"
	"black-shop/internal/adapters/storage/postgresql"
	"black-shop/internal/app/identity"
	"black-shop/pkg/config"
	"black-shop/pkg/logger"
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
		appLogger.Fatal("identity-service failed to run", zap.Error(err))
	}
}

func run(ctx context.Context, cfg *config.Config, appLogger *zap.Logger) error {
	// --- Database Connection ---
	appLogger.Info("connecting to PostgreSQL...")
	dbConn, err := pgx.Connect(ctx, cfg.PostgresIdentityURI)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer func(dbConn *pgx.Conn, ctx context.Context) {
		_ = dbConn.Close(ctx)
	}(dbConn, ctx)
	appLogger.Info("Successfully connected to PostgreSQL")

	// --- Layers ---
	userRepo, err := postgresql.NewUserRepository(dbConn, appLogger)
	if err != nil {
		return fmt.Errorf("failed to create user repository: %w", err)
	}

	// --- JWT ---
	jwtSecretKey := "a_very_secret_key"
	tokenManager := auth.NewTokenManager(jwtSecretKey)

	identityService := identity.NewService(userRepo, appLogger, tokenManager)
	grpcHandler := identity.NewGRPCServer(identityService, appLogger)

	// --- Servers ---
	errCh := make(chan error, 1)
	go func() {
		errCh <- runGRPCServer(cfg.IdentityGRPCPort, grpcHandler, tokenManager, appLogger)
	}()
	go func() {
		errCh <- runRESTGateway(ctx, cfg.IdentityHTTPPort, cfg.IdentityGRPCAddr, appLogger)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("a server failed: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func runGRPCServer(port string, handler pb.IdentityServiceServer, tokenManager *auth.TokenManager, appLogger *zap.Logger) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", port, err)
	}

	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(tokenManager.AuthenticationInterceptor))
	pb.RegisterIdentityServiceServer(gRPCServer, handler)

	appLogger.Info("gRPC Server is running", zap.String("port", port))
	return gRPCServer.Serve(lis)
}

func runRESTGateway(ctx context.Context, httpPort, grpcAddr string, appLogger *zap.Logger) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterIdentityServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	handler := cors.AllowAll().Handler(mux)

	appLogger.Info("HTTP REST Gateway is running", zap.String("port", httpPort))
	return http.ListenAndServe(httpPort, handler)
}
