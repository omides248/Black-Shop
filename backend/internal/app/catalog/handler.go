package catalog

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "black-shop-service/api/proto/v1"
	"black-shop-service/internal/domain/catalog"
)

type GRPCServer struct {
	pb.UnimplementedCatalogServiceServer
	service *Service
	logger  *zap.Logger
}

func NewGRPCServer(service *Service, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		service: service,
		logger:  logger.Named("catalog_grpc_handler"),
	}
}

func (s *GRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	s.logger.Info("received GetProduct request", zap.String("product_id", req.GetId()))

	productID := catalog.ProductID(req.GetId())
	foundProduct, err := s.service.GetProduct(ctx, productID)
	if err != nil {
		s.logger.Error("failed to get product from service",
			zap.String("product_id", req.GetId()),
			zap.Error(err),
		)

		if errors.Is(err, catalog.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, "product with id '%s' not found", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "internal server error while getting product")
	}

	return &pb.Product{
		Id:   string(foundProduct.ID),
		Name: foundProduct.Name,
	}, nil
}

func (s *GRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	s.logger.Info("received ListProducts request")

	products, err := s.service.FindAllProducts(ctx)
	if err != nil {
		s.logger.Error("failed to get all products from service", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to retrieve product list")
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:   string(p.ID),
			Name: p.Name,
		})
	}

	s.logger.Info("successfully retrieved all products", zap.Int("count", len(pbProducts)))
	return &pb.ListProductsResponse{Products: pbProducts}, nil
}

func (s *GRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	s.logger.Info("received CreateProduct request", zap.String("name", req.GetName()))

	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "product name cannot be empty")
	}

	product, err := s.service.CreateProduct(ctx, req.GetName())
	if err != nil {
		s.logger.Error("failed to create product via service", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create product")
	}

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:   string(product.ID),
			Name: product.Name,
		},
	}, nil
}
