package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"vinted/otel-workshop/internal/config"
	"vinted/otel-workshop/internal/factory"
	"vinted/otel-workshop/internal/telemetry"

	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"
)

const (
	name = "vinted/otel-workshop/factory"
)

var (
	tracer = otel.Tracer(name)
)

type FactoryConfig struct {
	FactoryAddress          string        `envconfig:"FACTORY_SERVICE_ADDR" validate:"required"`
	KafkaBrokers            []string      `envconfig:"KAFKA_SERVICE_ADDR" validate:"required"`
	FactoryKafkaTopic       string        `envconfig:"FACTORY_SERVICE_KAFKA_TOPIC" validate:"required"`
	FactoryMaxProduction    int           `envconfig:"FACTORY_SERVICE_MAX_PRODUCTION" validate:"required"`
	FactoryShippingInterval time.Duration `envconfig:"FACTORY_SERVICE_SHIPPING_INTERVAL" validate:"required"`
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

	cfg, err := config.Load[FactoryConfig]()
	if err != nil {
		logger.Error("new config", "error", err)
		os.Exit(1)
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		shipper, err := factory.NewKafkaShipper(logger, cfg.KafkaBrokers, cfg.FactoryKafkaTopic)
		if err != nil {
			logger.Error("failed to create Kafka shipper", "error", err)
			return err
		}

		factory := factory.NewProductFactory(logger, cfg.FactoryMaxProduction, shipper)

		ticker := time.NewTicker(cfg.FactoryShippingInterval)
		defer ticker.Stop()

		for range ticker.C {
			ctx, span := tracer.Start(ctx, "factory.Produce")

			err = factory.Produce(ctx)
			if err != nil {
				logger.Error("failed to produce products", "error", err)
				return err
			}

			span.End()
		}

		return nil
	})

	g.Go(func() error {
		orderShipper, err := factory.NewKafkaShipper(logger, cfg.KafkaBrokers, cfg.FactoryKafkaTopic)
		if err != nil {
			logger.Error("failed to create orders Kafka shipper", "error", err)
			return err
		}

		server := factory.NewFactoryServer(logger, cfg.FactoryAddress, orderShipper)

		return server.StartAndRun()
	})

	if err := g.Wait(); err != nil {
		logger.Error("factory failed", "error", err)
		os.Exit(1)
	}
}
