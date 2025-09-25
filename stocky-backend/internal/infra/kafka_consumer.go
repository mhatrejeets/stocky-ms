package infra

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type RewardConsumer struct{}

func (rc *RewardConsumer) Start(ctx context.Context, brokers []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	client, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		logrus.WithError(err).Error("Failed to start Kafka consumer")
		return
	}
	defer client.Close()
	partitionConsumer, err := client.ConsumePartition("reward-events", 0, sarama.OffsetNewest)
	if err != nil {
		logrus.WithError(err).Error("Failed to consume partition")
		return
	}
	defer partitionConsumer.Close()
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			logrus.Infof("Received event: %s", string(msg.Value))
		case <-ctx.Done():
			return
		}
	}
}
