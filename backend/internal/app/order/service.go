package order

import (
	"black-shop-service/internal/domain/order"
	"context"
	"errors"
	"go.uber.org/zap"
)

type Service struct {
	cartRepo  order.CartRepository
	orderRepo order.OrderRepository
	logger    *zap.Logger
}

func NewService(cartRepo order.CartRepository, orderRepo order.OrderRepository, logger *zap.Logger) *Service {
	return &Service{
		cartRepo:  cartRepo,
		orderRepo: orderRepo,
		logger:    logger.Named("order_service"),
	}
}

func (s *Service) AddItemToCart(ctx context.Context, userID string, item order.CartItem) (*order.Cart, error) {
	s.logger.Info("adding item to cart", zap.String("user_id", userID), zap.String("product_id", item.ProductID))

	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, order.ErrCartNotFound) {
		s.logger.Error("failed to get existing cart", zap.Error(err))
		return nil, err
	}

	if cart == nil {
		cart = &order.Cart{UserID: userID, Items: []order.CartItem{}}
	}

	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, item)
	}

	if err := s.cartRepo.Save(ctx, cart); err != nil {
		s.logger.Error("failed to save updated cart", zap.Error(err))
		return nil, err
	}

	return cart, nil
}

func (s *Service) GetCart(ctx context.Context, userID string) (*order.Cart, error) {
	s.logger.Info("getting user cart", zap.String("user_id", userID))
	return s.cartRepo.GetByUserID(ctx, userID)
}

func (s *Service) CreateOrderFromCart(ctx context.Context, userID string) (*order.Order, error) {
	s.logger.Info("creating order from cart", zap.String("user_id", userID))

	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cannot create order from an empty cart")
	}

	var orderItems []order.OrderItem
	for _, cartItem := range cart.Items {

		// TODO Request to GRPC for get product price of catalog-service
		price := 10.00

		orderItems = append(orderItems, order.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     price,
		})
	}

	newOrder, err := order.NewOrder(userID, orderItems)
	if err != nil {
		return nil, err
	}

	if err := s.orderRepo.Save(ctx, newOrder); err != nil {
		return nil, err
	}

	if err := s.cartRepo.Delete(ctx, userID); err != nil {
		s.logger.Error("failed to delete cart after order creation", zap.String("user_id", userID), zap.Error(err))
	}

	s.logger.Info("order created successfully", zap.String("order_id", string(newOrder.ID)))
	return newOrder, nil

}
