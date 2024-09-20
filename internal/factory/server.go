package factory

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"vinted/otel-workshop/pb/genproto/otelworkshop"
)

type FactoryServer struct {
	logger         *slog.Logger
	shipper        Shipper
	factoryAddress string
}

func NewFactoryServer(logger *slog.Logger, factoryAddress string, shipper Shipper) *FactoryServer {
	return &FactoryServer{
		logger:         logger,
		shipper:        shipper,
		factoryAddress: factoryAddress,
	}
}

func (s *FactoryServer) StartAndRun() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/make", s.handleMake)

	err := http.ListenAndServe(s.factoryAddress, mux)
	if err != nil {
		s.logger.Error("failed to serve HTTP", "error", err)
		return err
	}

	return nil
}

func (s *FactoryServer) handleMake(w http.ResponseWriter, r *http.Request) {
	var p otelworkshop.Product

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.logger.Info("received order to make", "name", p.Name, "color", p.Color, "quantity", p.Quantity)

	var products []*otelworkshop.Product

	for i := 0; i < int(p.Quantity); i++ {
		products = append(products, &otelworkshop.Product{
			Name:     p.Name,
			Color:    p.Color,
			Quantity: 1,
		})
	}

	err = s.shipper.Ship(context.Background(), products)
	if err != nil {
		s.logger.Error("failed to ship product", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
