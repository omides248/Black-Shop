package main

import (
	"catalog/config"
	"catalog/internal/adapters"
	"catalog/internal/application"
	"catalog/internal/delivery/http/error_mapping"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"pkg/echo/error_handler"
	"pkg/logger"
	"pkg/minio"
	"pkg/validation"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "catalog/api/proto/v1"

	grpcserver "catalog/internal/delivery/grpc"
	httpserver "catalog/internal/delivery/http/router"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger.Init("development")
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(logger.Logger)

	appLogger := logger.Logger
	defer func(appLogger *zap.Logger) {
		_ = appLogger.Sync()
	}(appLogger)

	if err := run(context.Background(), cfg, appLogger); err != nil {
		appLogger.Fatal("catalog-service failed to run", zap.Error(err))
	}
}

func run(ctx context.Context, cfg *config.Config, appLogger *zap.Logger) error {
	// Setup database
	db, err := setupDatabase(ctx, cfg.Database.MongoURI, appLogger)
	if err != nil {
		appLogger.Error("failed to setup database", zap.Error(err))
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(db.Client(), ctx)

	// Setup OpenTelemetry Providers (Tracer and Meter)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//tp, mp, err := newOtelProviders(ctx, appLogger)
	//if err != nil {
	//	appLogger.Fatal("failed to create otel providers", zap.Error(err))
	//}
	//defer func() {
	//	if err := tp.Shutdown(ctx); err != nil {
	//		appLogger.Fatal("failed to shutdown TracerProvider", zap.Error(err))
	//	}
	//	if err := mp.Shutdown(ctx); err != nil {
	//		appLogger.Fatal("failed to shutdown MeterProvider", zap.Error(err))
	//	}
	//}()
	appLogger.Info("OpenTelemetry providers initialized")

	// --- Adapters ---
	adapter, err := adapters.NewAdapter(db, appLogger)
	if err != nil {
		appLogger.Fatal("failed to create adapters", zap.Error(err))
	}

	// --- MinIO Service ---
	appLogger.Info("connecting to MinIO...")
	minioService, err := minio.NewService(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.PublicURL,
		appLogger,
	)
	if err != nil {
		return fmt.Errorf("failed to create MinIO service: %w", err)
	}

	// --- Local Storage ---
	//localStorageService, err := local_storage.NewService(
	//	cfg.LocalStorage.PublicStoragePath,
	//	appLogger,
	//)
	//if err != nil {
	//	return fmt.Errorf("failed to create local storage service: %w", err)
	//}

	// --- Application Services ---
	service := application.NewService(adapter.ProductRepo, adapter.CategoryRepo, appLogger)

	grpcServerDeps := grpcserver.ServerDependencies{
		ProductService:  service.ProductService,
		CategoryService: service.CategoryService,
		Logger:          appLogger,
	}

	// --- Delivery Layer (gRPC Handler) ---
	grpcServerHandler := grpcserver.NewServer(grpcServerDeps, appLogger)

	// Setup gRPC
	errCh := make(chan error, 1)
	go func() {
		errCh <- runGRPCServer(cfg.General.GRPCPort, grpcServerHandler, appLogger)
	}()

	// Setup Echo
	go func() {
		errCh <- runHTTPServer(cfg.General.HTTPPort, service.CategoryService, service.ProductService, minioService, cfg, appLogger)
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
	appLogger.Info("starting gRPC server...", zap.String("port", port))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", port, err)
	}
	gRPCServer := grpc.NewServer()
	pb.RegisterCatalogServiceServer(gRPCServer, handler)
	appLogger.Info("gRPC Server is running on", zap.String("port", port))
	return gRPCServer.Serve(lis)
}

func runHTTPServer(port string, catSvc application.CategoryService, prodSvc application.ProductService, minioService *minio.Service, cfg *config.Config, appLogger *zap.Logger) error {
	appLogger.Info("starting HTTP (Echo) server...", zap.String("port", port))
	e := echo.New()

	e.Validator = validation.New()

	//e.Use(otelecho.Middleware("catalog-service"))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.RequestID())

	//e.Use(middleware.MetricsMiddleware())

	domainErrorMappings := error_mapping.GetDomainErrorMappings()

	e.HTTPErrorHandler = error_handler.NewHTTPErrorHandler(domainErrorMappings, appLogger)

	httpserver.Setup(e, catSvc, prodSvc, minioService, cfg, appLogger)

	appLogger.Info("HTTP (Echo) Server is running on", zap.String("port", port))
	if err := e.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("echo server failed: %w", err)
	}
	return nil
}

func newOtelProviders(ctx context.Context, appLogger *zap.Logger) (*sdktrace.TracerProvider, *metric.MeterProvider, error) {
	// 1. Create a gRPC trace exporter
	traceExp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("otel-collector:4317"))
	if err != nil {
		appLogger.Error("failed to create otlptracegrpc exporter", zap.Error(err))
		return nil, nil, err
	}

	// 2. Create a gRPC metric exporter
	metricExp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint("otel-collector:4317"))
	if err != nil {
		appLogger.Error("failed to create otlpmetricgrpc exporter", zap.Error(err))
		return nil, nil, err
	}

	// 3. Create a resource with service name
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName("catalog-service"),
	))
	if err != nil {
		appLogger.Error("failed to create resource", zap.Error(err))
		return nil, nil, err
	}

	// 4. Create a tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(res),
	)

	// 5. Create a meter provider
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExp, metric.WithInterval(5*time.Second))), // Push every 5 seconds
	)

	// 6. Set the global providers
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)

	return tp, mp, nil
}

func newMeterProvider(ctx context.Context, appLogger *zap.Logger) (*metric.MeterProvider, error) {
	// 1. Create a gRPC exporter
	exp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint("otel-collector:4317"))
	if err != nil {
		appLogger.Error("failed to create otlpmetricgrpc exporter", zap.Error(err))
		return nil, err
	}

	// 2. Create a resource with service name
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName("catalog-service"),
	))
	if err != nil {
		appLogger.Error("failed to create resource", zap.Error(err))
		return nil, err
	}

	// 3. Create a meter provider
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exp, metric.WithInterval(5*time.Second))), // Push every 5 seconds
	)

	// 4. Set the global meter provider
	otel.SetMeterProvider(mp)
	return mp, nil
}
