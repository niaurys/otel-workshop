package buyer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"vinted/otel-workshop/pb/genproto/otelworkshop"

	"github.com/sirupsen/logrus"
)

type BuyerServer struct {
	logger      *logrus.Logger
	factoryAddr string
	client      http.Client
}

func NewBuyerServer(logger *logrus.Logger, factoryAddr string, client http.Client) *BuyerServer {
	return &BuyerServer{
		logger:      logger,
		factoryAddr: factoryAddr,
		client:      client,
	}
}

func (s *BuyerServer) HandleOrder(w http.ResponseWriter, r *http.Request) {
	var p otelworkshop.Product

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.logger.WithFields(logrus.Fields{
		"name":     p.Name,
		"color":    p.Color,
		"quantity": p.Quantity,
	}).Info("received order")

	url := fmt.Sprintf("http://%s/make", s.factoryAddr)

	order, err := json.Marshal(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewBuffer(order))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	proxyReq.Header.Set("Host", r.Host)
	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	resp, err := s.client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
}
