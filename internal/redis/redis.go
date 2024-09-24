package redis

import (
	"context"
	"vinted/otel-workshop/pb/genproto/otelworkshop"

	redisotel "github.com/redis/go-redis/extra/redisotel/v9"
	redis "github.com/redis/go-redis/v9"
)

type RedisClient interface {
	DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd
	IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}

type WorkshopClient struct {
	client RedisClient
}

func NewWorkshopRedisClient(redisAddr string) *WorkshopClient {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil
	}

	return &WorkshopClient{
		client: client,
	}
}

func key(product *otelworkshop.Product) string {
	return product.Name + ":" + product.Color
}

func (r *WorkshopClient) Decrement(ctx context.Context, product *otelworkshop.Product, decrement int64) error {
	return r.client.DecrBy(ctx, key(product), decrement).Err()
}

func (r *WorkshopClient) Increment(ctx context.Context, product *otelworkshop.Product, value int64) error {
	return r.client.IncrBy(ctx, key(product), value).Err()
}

func (r *WorkshopClient) GetValue(ctx context.Context, product *otelworkshop.Product) (int64, error) {
	value, err := r.client.Get(ctx, key(product)).Int64()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return value, nil
}
