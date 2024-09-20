package main

import (
	"context"
	"log/slog"
	"os"

	"vinted/otel-workshop/internal/config"
	"vinted/otel-workshop/internal/warehouse"
)

type WarehouseConfig struct {
	KafkaBrokers           []string `envconfig:"KAFKA_SERVICE_ADDR" validate:"required"`
	RedisAddress           string   `envconfig:"REDIS_SERVICE_ADDR" validate:"required"`
	WarehouseTopic         string   `envconfig:"FACTORY_SERVICE_KAFKA_TOPIC" validate:"required"`
	WarehouseConsumerGroup string   `envconfig:"WAREHOUSE_SERVICE_CONSUMER_GROUP" validate:"required"`
}

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	cfg, err := config.Load[WarehouseConfig]()
	if err != nil {
		logger.Error("new config", "error", err)
		os.Exit(1)
	}

	storage := warehouse.NewRedisWarehouseStorage(logger, cfg.RedisAddress)

	warehouse, err := warehouse.NewKafkaRedisWarehouse(
		logger,
		cfg.KafkaBrokers,
		[]string{cfg.WarehouseTopic},
		cfg.WarehouseTopic,
		storage,
	)
	if err != nil {
		logger.Error("failed to create warehouse", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		if err := warehouse.PickAndStore(ctx); err != nil {
			logger.Error("failed to pick and store products", "error", err)
			os.Exit(1)
		}
	}
}
