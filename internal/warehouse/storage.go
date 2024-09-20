package warehouse

import (
	"context"
	"encoding/json"
	"log/slog"
	"vinted/otel-workshop/internal/redis"
	"vinted/otel-workshop/pb/genproto/otelworkshop"
)

type WarehouseStorage interface {
	Store(ctx context.Context, data []byte) error
}

type RedisWarehouseStorage struct {
	redisClient *redis.WorkshopClient
	logger      *slog.Logger
}

func NewRedisWarehouseStorage(logger *slog.Logger, addr string) *RedisWarehouseStorage {
	return &RedisWarehouseStorage{
		redisClient: redis.NewWorkshopRedisClient(addr),
		logger:      logger,
	}
}

func (s *RedisWarehouseStorage) Store(ctx context.Context, data []byte) error {
	var product otelworkshop.Product

	s.logger.Info("storing product", "data", string(data))

	err := json.Unmarshal(data, &product)
	if err != nil {
		s.logger.Error("failed to unmarshal message", "error", err)
		return err
	}

	err = s.redisClient.Increment(ctx, &product, 1)
	if err != nil {
		s.logger.Error("failed to store product", "error", err)
		return err
	}

	return nil
}
