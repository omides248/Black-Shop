package redis

import (
	"black-shop/internal/order/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type cartRepo struct {
	client *redis.Client
	logger *zap.Logger
}

func NewCartRepository(client *redis.Client, logger *zap.Logger) domain.CartRepository {
	return &cartRepo{
		client: client,
		logger: logger.Named("redis_cart_repo"),
	}
}

func generateCartKey(userID string) string {
	return fmt.Sprintf("cart:%s", userID)
}

func (r *cartRepo) Save(ctx context.Context, cart *domain.Cart) error {
	key := generateCartKey(cart.UserID)
	r.logger.Info("saving cart", zap.String("key", key))

	cartData, err := json.Marshal(cart.Items)
	if err != nil {
		r.logger.Error("failed to marshal cart data", zap.Error(err))
		return err
	}

	err = r.client.Set(ctx, key, cartData, time.Hour*24*7).Err()
	if err != nil {
		r.logger.Error("failed to save cart to redis", zap.Error(err))
		return err
	}

	return nil
}

func (r *cartRepo) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	key := generateCartKey(userID)
	r.logger.Info("getting cart", zap.String("key", key))

	cartData, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Warn("cart not found", zap.String("key", key))
			return nil, domain.ErrCartNotFound
		}
		r.logger.Error("failed to get cart from redis", zap.Error(err))
		return nil, err
	}

	var items []domain.CartItem
	if err := json.Unmarshal(cartData, &items); err != nil {
		r.logger.Error("failed to unmarshal cart data", zap.Error(err))
		return nil, err
	}

	return &domain.Cart{
		UserID: userID,
		Items:  items,
	}, nil
}

func (r *cartRepo) Delete(ctx context.Context, userID string) error {
	key := generateCartKey(userID)
	r.logger.Info("deleting cart", zap.String("key", key))

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("failed to delete cart from redis", zap.Error(err))
		return err
	}
	return nil
}
