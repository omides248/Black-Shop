package services

import (
	"context"
	"errors"
	"fmt"
	"order/internal/domain"
	"pkg/wallet"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	cartRepo      domain.CartRepository
	orderRepo     domain.OrderRepository
	walletService *wallet.Service
	logger        *zap.Logger
}

func NewService(
	cartRepo domain.CartRepository,
	orderRepo domain.OrderRepository,
	walletService *wallet.Service,
	logger *zap.Logger,
) *Service {

	return &Service{
		cartRepo:      cartRepo,
		orderRepo:     orderRepo,
		walletService: walletService,
		logger:        logger.Named("order_service"),
	}
}

func (s *Service) AddItemToCart(ctx context.Context, userID string, item domain.CartItem) (*domain.Cart, error) {
	s.logger.Info("adding item to cart", zap.String("user_id", userID), zap.String("product_id", item.ProductID))

	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, domain.ErrCartNotFound) {
		s.logger.Error("failed to get existing cart", zap.Error(err))
		return nil, err
	}

	if cart == nil {
		cart = &domain.Cart{UserID: userID, Items: []domain.CartItem{}}
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

	fmt.Println("*********************** Cart", cart)

	if err := s.cartRepo.Save(ctx, cart); err != nil {
		s.logger.Error("failed to save updated cart", zap.Error(err))
		return nil, err
	}

	return cart, nil
}

func (s *Service) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	s.logger.Info("getting user cart", zap.String("user_id", userID))
	return s.cartRepo.GetByUserID(ctx, userID)
}

func (s *Service) CreateOrderFromCart(ctx context.Context, userID string, paymentMethod string) (*domain.Order, error) {
	s.logger.Info("creating order from cart", zap.String("user_id", userID), zap.String("payment_method", paymentMethod))

	cart, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cannot create order from an empty cart")
	}

	var orderItems []domain.OrderItem
	for _, cartItem := range cart.Items {

		// TODO Request to GRPC for get product price of catalog-service
		price := 10.00

		orderItems = append(orderItems, domain.OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     price,
		})
	}

	newOrder, err := domain.NewOrder(userID, paymentMethod, orderItems)
	if err != nil {
		return nil, err
	}

	if paymentMethod == "CRYPTO" {
		s.logger.Info("crypto payment method selected, generating address...")

		// derivationPath (Mnemonic): m/purpose'/coin_type'/account'/change/address_index
		// purpose: Always is 44 for BIP44 standard
		// coin_type: Bitcoin = 0, Ethereum = 60
		// account: Support multi accounts ( personal or work account)
		// address_index: is a counter increase for a new address
		addressIndex := time.Now().Unix()
		derivationPath := fmt.Sprintf("m/44'/60'/0'/%d", addressIndex)

		paymentAddr, err := s.walletService.DeriveAddress(derivationPath)
		if err != nil {
			s.logger.Error("failed to derive payment address", zap.Error(err))
			return nil, fmt.Errorf("could not generate payment address for order")
		}

		newOrder.PaymentAddress = &paymentAddr
		newOrder.Status = domain.StatusAwaitingPayment
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
