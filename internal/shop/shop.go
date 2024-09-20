package shop

import (
	"context"
	"sync"
	"vinted/otel-workshop/internal/product"
	"vinted/otel-workshop/internal/redis"
	"vinted/otel-workshop/pb/genproto/otelworkshop"

	"go.uber.org/zap"
)

type RedisShop struct {
	redisClient *redis.WorkshopClient
	mux         sync.RWMutex
	inventory   []*otelworkshop.Product
	logger      *zap.Logger

	otelworkshop.UnimplementedShopServiceServer
}

func NewRedisShop(logger *zap.Logger, redisAddr string) *RedisShop {
	return &RedisShop{
		redisClient: redis.NewWorkshopRedisClient(redisAddr),
		logger:      logger,
	}
}

func (s *RedisShop) ListProducts(ctx context.Context, req *otelworkshop.Empty) (*otelworkshop.ListProductsResponse, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	s.logger.Info("listing products", zap.Int("count", len(s.inventory)))

	return &otelworkshop.ListProductsResponse{Products: s.inventory}, nil
}

func (s *RedisShop) BuyProduct(ctx context.Context, req *otelworkshop.BuyProductRequest) (*otelworkshop.Product, error) {
	s.logger.Info("buying product", zap.String("name", req.Name), zap.String("surname", req.Surname), zap.Any("product", req.Product))

	err := s.redisClient.Decrement(ctx, req.Product, req.Product.Quantity)
	if err != nil {
		s.logger.Error("failed to decrement product quantity", zap.Error(err))
		return nil, err
	}

	s.logger.Info("product bought", zap.String("name", req.Name), zap.String("surname", req.Surname), zap.Any("product", req.Product))

	return req.Product, err
}

func (s *RedisShop) UpdateInventory(ctx context.Context) error {
	inventory := make([]*otelworkshop.Product, 0)

	for _, color := range product.Colors() {
		for _, name := range product.Names() {
			product := &otelworkshop.Product{
				Name:  name,
				Color: color,
			}

			quantity, err := s.redisClient.GetValue(ctx, product)
			if err != nil {
				return err
			}

			product.Quantity = quantity

			inventory = append(inventory, product)
		}
	}

	s.mux.Lock()
	s.inventory = inventory
	s.mux.Unlock()

	return nil
}
