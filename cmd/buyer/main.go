package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"vinted/otel-workshop/internal/buyer"
	"vinted/otel-workshop/internal/config"
	"vinted/otel-workshop/internal/telemetry"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/sync/errgroup"
)

const (
	name = "vinted/otel-workshop/buyer"
)

var (
	tracer = otel.Tracer(name)
	meter  = otel.Meter(name)
)

type BuyerConfig struct {
	BuyerAddress   string        `envconfig:"BUYER_SERVICE_ADDR" validate:"required"`
	BuyingInterval time.Duration `envconfig:"BUYER_SERVICE_BUY_INTERVAL" validate:"required"`
	ShopAddress    string        `envconfig:"SHOP_SERVICE_ADDR" validate:"required"`
	FactoryAddress string        `envconfig:"FACTORY_SERVICE_ADDR" validate:"required"`
}

func main() {
	ctx := context.Background()

	shutdown, err := telemetry.SetupOtelSDK(ctx)
	defer func() {
		if err := shutdown(ctx); err != nil {
			logrus.Fatalf("failed to shutdown otel sdk: %v", err)
		}
	}()
	if err != nil {
		logrus.Fatalf("failed to setup otel sdk: %v", err)
	}

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Load[BuyerConfig]()
	if err != nil {
		logger.Fatalf("new config: %v", err)
	}
	logger.WithFields(logrus.Fields{
		"buyer_address":   cfg.BuyerAddress,
		"buying_interval": cfg.BuyingInterval,
		"shop_address":    cfg.ShopAddress,
		"factory_address": cfg.FactoryAddress,
	}).Info("starting buyer service")

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		client := http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}

		server := buyer.NewBuyerServer(logger, cfg.FactoryAddress, client)

		mux := http.NewServeMux()
		mux.HandleFunc("/order", server.HandleOrder)
		mux.Handle("/metrics", promhttp.Handler())

		err = http.ListenAndServe(cfg.BuyerAddress, mux)
		if err != nil {
			logger.Fatalf("server failed: %v", err)
			return err
		}

		return nil
	})

	orders, err := meter.Int64Counter(
		"buyer.orders",
		metric.WithDescription("Number of orders made"),
		metric.WithUnit("count"),
	)
	if err != nil {
		logger.Fatalf("failed to create orders counter: %v", err)
	}

	buyDuration, err := meter.Int64Histogram(
		"buyer.buy_duration",
		metric.WithDescription("Duration of buying products"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(0, 100, 200, 300, 400, 500, 600, 700),
	)
	if err != nil {
		logger.Fatalf("failed to create buy duration histogram: %v", err)
	}

	g.Go(func() error {
		buyer, err := buyer.NewRandomBuyer(logger, cfg.ShopAddress)
		if err != nil {
			logger.Fatalf("failed to create buyer: %v", err)
		}

		ticker := time.NewTicker(cfg.BuyingInterval)
		defer ticker.Stop()

		for range ticker.C {
			id := uuid.New()

			memeber, err := baggage.NewMember("order_id", id.String())
			if err != nil {
				logger.Fatalf("failed to create baggage member: %v", err)
			}

			bag, err := baggage.New(memeber)
			if err != nil {
				logger.Fatalf("failed to create baggage: %v", err)
			}

			ctx := baggage.ContextWithBaggage(ctx, bag)

			ctx, span := tracer.Start(ctx, "buyer.Buy")
			orders.Add(ctx, 1)

			logger.Info("buying product")
			now := time.Now()
			err = buyer.Buy(ctx)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to buy")
				logger.Fatalf("failed to buy: %v", err)
				return err
			}
			buyDuration.Record(ctx, time.Since(now).Milliseconds())

			span.End()
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Fatalf("buyer failed: %v", err)
	}

}
