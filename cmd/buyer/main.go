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
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
)

const (
	name = "vinted/otel-workshop/buyer"
)

var (
	tracer = otel.Tracer(name)
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

		err = http.ListenAndServe(cfg.BuyerAddress, mux)
		if err != nil {
			logger.Fatalf("server failed: %v", err)
			return err
		}

		return nil
	})

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

			logger.Info("buying product")
			if err := buyer.Buy(ctx); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to buy")
				logger.Fatalf("failed to buy: %v", err)
				return err
			}

			span.End()
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Fatalf("buyer failed: %v", err)
	}

}
