package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vinted/otel-workshop/internal/config"
	"vinted/otel-workshop/internal/shop"
	"vinted/otel-workshop/pb/genproto/otelworkshop"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ShopConfig struct {
	RedisAddress                string        `envconfig:"REDIS_SERVICE_ADDR" validate:"required"`
	ShopAddress                 string        `envconfig:"SHOP_SERVICE_ADDR" validate:"required"`
	ShopInventoryUpdateInterval time.Duration `envconfig:"SHOP_SERVICE_INVENTORY_UPDATE_INTERVAL" validate:"required"`
}

func main() {
	logger := zap.Must(zap.NewProduction())
	defer func() {
		_ = logger.Sync()
	}()

	cfg, err := config.Load[ShopConfig]()
	if err != nil {
		logger.Fatal("new config", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shop := shop.NewRedisShop(logger, cfg.RedisAddress)
	if err = shop.UpdateInventory(ctx); err != nil {
		logger.Fatal("failed to update inventory", zap.Error(err))
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		ticker := time.NewTicker(cfg.ShopInventoryUpdateInterval)
		defer ticker.Stop()

		for range ticker.C {
			logger.Info("updating inventory")
			if err := shop.UpdateInventory(ctx); err != nil {
				logger.Error("failed to update inventory", zap.Error(err))
				return err
			}
		}

		return nil
	})

	g.Go(func() error {
		listen, err := net.Listen("tcp", cfg.ShopAddress)
		if err != nil {
			logger.Error("failed to listen", zap.String("address", cfg.ShopAddress), zap.Error(err))
			return err
		}

		grpcServer := grpc.NewServer()
		otelworkshop.RegisterShopServiceServer(grpcServer, shop)
		reflection.Register(grpcServer)

		logger.Info("starting server", zap.String("address", cfg.ShopAddress))
		if err := grpcServer.Serve(listen); err != nil {
			logger.Error("failed to serve", zap.Error(err))
			return err
		}

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		grpcServer.GracefulStop()

		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Fatal("shop failed", zap.Error(err))
	}
}
