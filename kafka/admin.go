package kafka

import (
	"context"
	"errors"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type TopicConfig struct {
	Name              string
	NumPartitions     int
	ReplicationFactor int
}

func EnsureTopics(ctx context.Context, topics ...TopicConfig) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return ErrEnsureTopics.WithCause(err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return ErrEnsureTopics.WithCause(err)
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))

	controllerConn, err := kafka.DialContext(ctx, "tcp", controllerAddr)
	if err != nil {
		return ErrEnsureTopics.WithCause(err)
	}
	defer controllerConn.Close()

	for _, t := range topics {
		err := controllerConn.CreateTopics(kafka.TopicConfig{
			Topic:             t.Name,
			NumPartitions:     t.NumPartitions,
			ReplicationFactor: t.ReplicationFactor,
		})
		if err != nil && !isTopicAlreadyExists(err) {
			return ErrEnsureTopics.WithCause(err)
		}
	}

	return nil
}

func isTopicAlreadyExists(err error) bool {
	var kafkaErr kafka.Error
	return errors.As(err, &kafkaErr) && kafkaErr == kafka.TopicAlreadyExists
}
