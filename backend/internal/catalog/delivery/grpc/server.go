package grpc

import (
	pb "black-shop/api/proto/catalog/v1"
	"black-shop/internal/catalog/application"
	"go.uber.org/zap"
)

//var _ pb.CatalogServiceServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedCatalogServiceServer
	service *application.Service
	logger  *zap.Logger
}

func NewServer(service *application.Service, logger *zap.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger.Named("catalog_grpc_handler"),
	}
}
