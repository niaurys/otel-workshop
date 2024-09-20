package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"vinted/otel-workshop/internal/buyer"
	"vinted/otel-workshop/internal/config"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BuyerConfig struct {
	BuyerAddress   string        `envconfig:"BUYER_SERVICE_ADDR" validate:"required"`
	BuyingInterval time.Duration `envconfig:"BUYER_SERVICE_BUY_INTERVAL" validate:"required"`
	ShopAddress    string        `envconfig:"SHOP_SERVICE_ADDR" validate:"required"`
	FactoryAddress string        `envconfig:"FACTORY_SERVICE_ADDR" validate:"required"`
}

func main() {
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

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		server := buyer.NewBuyerServer(logger, cfg.FactoryAddress, http.Client{})

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
			logger.Info("buying product")
			if err := buyer.Buy(ctx); err != nil {
				logger.Fatalf("failed to buy: %v", err)
				return err
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Fatalf("buyer failed: %v", err)
	}

}
