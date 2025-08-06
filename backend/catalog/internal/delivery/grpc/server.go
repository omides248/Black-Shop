package grpc

import (
	pb "catalog/api/proto/v1"
	"catalog/internal/application"
	"go.uber.org/zap"
)

//var _ pb.CatalogServiceServer = (*Server)(nil)

type ServerDependencies struct {
	ProductService  application.ProductService
	CategoryService application.CategoryService
	Logger          *zap.Logger
}

type Server struct {
	pb.UnimplementedCatalogServiceServer
	productService  application.ProductService
	categoryService application.CategoryService
	logger          *zap.Logger
}

func NewServer(deps ServerDependencies, logger *zap.Logger) *Server {
	return &Server{
		productService:  deps.ProductService,
		categoryService: deps.CategoryService,
		logger:          logger.Named("catalog_grpc_handler"),
	}
}
