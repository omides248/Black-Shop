package grpc

import (
	"black-shop/internal/catalog/domain"
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "black-shop/api/proto/catalog/v1"
)

func (s *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	s.logger.Info("received CreateCategory request", zap.String("name", req.GetName()))

	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "category name cannot be empty")
	}

	category, err := s.service.CreateCategory(ctx, req.GetName(), req.ImageUrl, req.ParentId)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return nil, status.Errorf(codes.NotFound, "parent category not found")
		}
		if errors.Is(err, domain.ErrCategoryDepthLimitExceeded) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		if errors.Is(err, domain.ErrCategoryHasProducts) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "failed to create category")
	}

	return &pb.CreateCategoryResponse{
		Category: &pb.Category{
			Id:       string(category.ID),
			Name:     category.Name,
			ImageUrl: category.ImageURL,
			ParentId: (*string)(category.ParentID),
			Depth:    int32(category.Depth),
		},
	}, nil
}

func (s *Server) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	s.logger.Info("received ListCategories request")

	categories, err := s.service.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error("failed to get all categories from service", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to retrieve category list")
	}

	// Convert domain's model to grpc's model
	pbCategories := make([]*pb.Category, len(categories))
	for i, c := range categories {
		pbCategories[i] = &pb.Category{
			Id:       string(c.ID),
			Name:     c.Name,
			ImageUrl: c.ImageURL,
			ParentId: (*string)(c.ParentID),
			Depth:    int32(c.Depth),
		}
	}

	return &pb.ListCategoriesResponse{
		Categories: pbCategories,
	}, nil
}
