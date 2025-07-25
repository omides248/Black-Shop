package order

import (
	pb "black-shop-service/api/proto/v1"
	"black-shop-service/internal/domain/order"
	"black-shop-service/pkg/contextkeys"
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedOrderServiceServer
	service *Service
	logger  *zap.Logger
}

func NewGRPCServer(server *Service, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		service: server,
		logger:  logger.Named("order_grpc_handler"),
	}
}

func (s *GRPCServer) AddItemToCart(ctx context.Context, req *pb.AddItemToCartRequest) (*pb.Cart, error) {
	userID, err := contextkeys.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	s.logger.Info("received AddItemToCart request", zap.String("user_id", userID), zap.String("product_id", req.ProductId))

	item := order.CartItem{
		ProductID: req.ProductId,
		Quantity:  int(req.Quantity),
	}

	cart, err := s.service.AddItemToCart(ctx, userID, item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add item to cart")
	}

	pbItems := make([]*pb.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		pbItems[i] = &pb.CartItem{ProductId: item.ProductID, Quantity: int32(item.Quantity)}
	}

	return &pb.Cart{UserId: cart.UserID, Items: pbItems}, nil
}

func (s *GRPCServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.Cart, error) {
	userID, err := contextkeys.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	s.logger.Info("received GetCart request", zap.String("user_id", userID))

	cart, err := s.service.GetCart(ctx, userID)
	if err != nil {
		if errors.Is(err, order.ErrCartNotFound) {
			return &pb.Cart{UserId: userID, Items: []*pb.CartItem{}}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to get cart")
	}

	pbItems := make([]*pb.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		pbItems[i] = &pb.CartItem{ProductId: item.ProductID, Quantity: int32(item.Quantity)}
	}

	return &pb.Cart{UserId: cart.UserID, Items: pbItems}, nil
}

func (s *GRPCServer) CreateOrderFromCart(ctx context.Context, req *pb.CreateOrderFromCartRequest) (*pb.Order, error) {
	userID, err := contextkeys.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	s.logger.Info("received CreateOrderFromCart request", zap.String("user_id", userID))

	newOrder, err := s.service.CreateOrderFromCart(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order")
	}

	pbOrderItems := make([]*pb.OrderItem, len(newOrder.Items))
	for i, item := range newOrder.Items {
		pbOrderItems[i] = &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
	}

	return &pb.Order{
		Id:         string(newOrder.ID),
		UserId:     newOrder.UserID,
		TotalPrice: newOrder.TotalPrice,
		Status:     string(newOrder.Status),
		Items:      pbOrderItems,
	}, nil
}
