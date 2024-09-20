package warehouse

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IBM/sarama"
)

type Warehouse interface {
	PickAndStore(ctx context.Context) error
}

type KafkaRedisWarehouse struct {
	consumerGroup sarama.ConsumerGroup
	handler       sarama.ConsumerGroupHandler
	topics        []string
	logger        *slog.Logger
}

func NewKafkaRedisWarehouse(logger *slog.Logger, brokerAddresses, topics []string, groupID string, storage WarehouseStorage) (*KafkaRedisWarehouse, error) {
	saramaConfig := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup(brokerAddresses, groupID, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaRedisWarehouse{
		consumerGroup: consumerGroup,
		handler: &productHandler{
			storage: storage,
			ready:   make(chan bool),
			logger:  logger,
		},
		topics: topics,
		logger: logger,
	}, nil
}

func (w *KafkaRedisWarehouse) PickAndStore(ctx context.Context) error {
	if err := w.consumerGroup.Consume(ctx, w.topics, w.handler); err != nil {
		if errors.Is(err, sarama.ErrClosedConsumerGroup) {
			w.logger.Info("consumer group closed")
			return err
		}
		w.logger.Error("failed to consume messages", "error", err)
	}

	if ctx.Err() != nil {
		w.logger.Error("context canceled", "error", ctx.Err())
		return ctx.Err()
	}

	return nil
}

type productHandler struct {
	storage WarehouseStorage
	ready   chan bool
	logger  *slog.Logger
}

func (p *productHandler) Setup(sarama.ConsumerGroupSession) error {
	close(p.ready)
	return nil
}

func (p *productHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (p *productHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				p.logger.Info("message channel was closed")
				return nil
			}

			p.logger.Info("message claimed", "value", string(message.Value), "timestamp", message.Timestamp, "topic", message.Topic)

			err := p.storage.Store(session.Context(), message.Value)
			if err != nil {
				p.logger.Error("failed to store", "error", err)
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
