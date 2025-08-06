package grpc

import (
	pb "catalog/api/proto/v1"
	"catalog/internal/domain"
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	s.logger.Info("received GetProduct request", zap.String("product_id", req.GetId()))

	productID := domain.ProductID(req.GetId())
	foundProduct, err := s.productService.GetProduct(ctx, productID)
	if err != nil {
		s.logger.Error("failed to get product from service",
			zap.String("product_id", req.GetId()),
			zap.Error(err),
		)

		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, "product with id '%s' not found", req.GetId())
		}
		return nil, status.Errorf(codes.Internal, "internal server error while getting product")
	}

	return &pb.Product{
		Id:   string(foundProduct.ID),
		Name: foundProduct.Name,
	}, nil
}

func (s *Server) ListProducts(ctx context.Context, _ *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	s.logger.Info("received ListProducts request")

	products, err := s.productService.FindAllProducts(ctx)
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

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	s.logger.Info("received CreateProduct request", zap.String("name", req.GetName()))

	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "product name cannot be empty")
	}

	product, err := s.productService.CreateProduct(ctx, req.GetName())
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
