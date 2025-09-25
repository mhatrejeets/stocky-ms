package infra

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
}

type RewardCreatedEvent struct {
	RewardID      string `json:"reward_id"`
	UserID        string `json:"user_id"`
	StockSymbol   string `json:"stock_symbol"`
	Shares        string `json:"shares"`
	RewardedAt    string `json:"rewarded_at"`
	CorrelationID string `json:"correlation_id"`
}

func (kp *KafkaProducer) PublishRewardCreated(ctx context.Context, event RewardCreatedEvent) error {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	correlationID := event.CorrelationID
	msg := &sarama.ProducerMessage{
		Topic: "reward-events",
		Value: sarama.ByteEncoder(msgBytes),
		Headers: []sarama.RecordHeader{{
			Key:   []byte("correlation_id"),
			Value: []byte(correlationID),
		}},
	}
	_, _, err = kp.Producer.SendMessage(msg)
	if err != nil {
		logrus.WithError(err).Error("Failed to publish reward event")
	}
	return err
}

// Noop fallback
func NoopPublishRewardCreated(ctx context.Context, event RewardCreatedEvent) error {
	logrus.Info("Noop publish: ", event)
	return nil
}
