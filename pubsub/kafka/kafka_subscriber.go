package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/purposeinplay/go-commons/pubsub"
	"go.uber.org/zap"
)

var _ pubsub.Subscriber[[]byte] = (*Subscriber)(nil)

// Subscriber represents a kafka subscriber.
type Subscriber struct {
	kafkaSubscriber *kafka.Subscriber
	clusterAdmin    sarama.ClusterAdmin
}

// NewSubscriber creates a new kafka subscriber.
func NewSubscriber(
	logger *zap.Logger,
	saramaConfig *sarama.Config,
	brokers []string,
	consumerGroup string,
) (*Subscriber, error) {
	saramaClient, err := sarama.NewClusterAdmin(brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("new sarama client: %w", err)
	}

	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaConfig,
			ConsumerGroup:         consumerGroup,
		},
		newLoggerAdapter(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("new kafka subscriber: %w", err)
	}

	return &Subscriber{
		kafkaSubscriber: sub,
		clusterAdmin:    saramaClient,
	}, nil
}

// Subscribe subscribes to a kafka topic.
func (s Subscriber) Subscribe(channels ...string) (pubsub.Subscription[[]byte], error) {
	if len(channels) != 1 {
		return nil, pubsub.ErrExactlyOneChannelAllowed
	}

	mes, err := s.kafkaSubscriber.Subscribe(context.Background(), channels[0])
	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	return newSubscription(mes, s.clusterAdmin), nil
}

// Close closes the kafka subscriber.
func (s Subscriber) Close() error {
	return s.kafkaSubscriber.Close()
}
