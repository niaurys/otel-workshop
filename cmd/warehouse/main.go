package main

import (
	"context"
	"log/slog"
	"os"

	"vinted/otel-workshop/internal/config"
	"vinted/otel-workshop/internal/telemetry"
	"vinted/otel-workshop/internal/warehouse"

	"go.opentelemetry.io/otel"
)

const (
	name = "vinted/otel-workshop/warehouse"
)

var (
	tracer = otel.Tracer(name)
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

	ctx := context.Background()

	shutdown, err := telemetry.SetupOtelSDK(ctx)
	defer func() {
		if err := shutdown(ctx); err != nil {
			logger.Error("failed to shutdown otel sdk", "error", err)
		}
	}()
	if err != nil {
		logger.Error("failed to setup otel sdk", "error", err)
		os.Exit(1)
	}

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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		ctx, span := tracer.Start(ctx, "warehouse.PickAndStore")

		if err := warehouse.PickAndStore(ctx); err != nil {
			logger.Error("failed to pick and store products", "error", err)
			os.Exit(1)
		}

		span.End()
	}
}
