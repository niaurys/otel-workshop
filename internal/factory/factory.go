package factory

import (
	"context"
	"encoding/json"
	"log/slog"
	"vinted/otel-workshop/internal/product"
	"vinted/otel-workshop/internal/random"
	"vinted/otel-workshop/pb/genproto/otelworkshop"

	"github.com/IBM/sarama"
)

type Shipper interface {
	Ship(context.Context, []*otelworkshop.Product) error
}

type Factory interface {
	Produce(context.Context) error
}

type ProductFactory struct {
	maxProduction int
	shipper       Shipper
	logger        *slog.Logger
}

func NewProductFactory(logger *slog.Logger, maxProduction int, shipper Shipper) *ProductFactory {
	return &ProductFactory{
		maxProduction: maxProduction,
		shipper:       shipper,
		logger:        logger,
	}
}

func (f *ProductFactory) Produce(ctx context.Context) error {
	var products []*otelworkshop.Product

	for i := 0; i < random.Int(f.maxProduction); i++ {
		products = append(products, product.New())
	}

	f.logger.Info("produced products", "count", len(products))

	return f.shipper.Ship(ctx, products)
}

type KafkaShipper struct {
	topic    string
	producer sarama.SyncProducer
	logger   *slog.Logger
}

func NewKafkaShipper(logger *slog.Logger, brokerAddresses []string, topic string) (*KafkaShipper, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerAddresses, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaShipper{
		topic:    topic,
		producer: producer,
		logger:   logger,
	}, nil
}

func (s *KafkaShipper) Ship(ctx context.Context, products []*otelworkshop.Product) error {
	var messages []*sarama.ProducerMessage

	for _, product := range products {
		productJson, err := json.Marshal(product)
		if err != nil {
			return err
		}

		messages = append(messages, &sarama.ProducerMessage{
			Topic: s.topic,
			Value: sarama.ByteEncoder(productJson),
		})
	}

	err := s.producer.SendMessages(messages)
	if err != nil {
		return err
	}

	s.logger.Info("shipped products", "count", len(products))

	return nil
}
