package infra

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/mhatrejeets/stocky-ms/internal/model"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
}


func (kp *KafkaProducer) PublishRewardCreated(ctx context.Context, event model.RewardCreatedEvent) error {
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
func NoopPublishRewardCreated(ctx context.Context, event model.RewardCreatedEvent) error {
	logrus.Info("Noop publish: ", event)
	return nil
}
