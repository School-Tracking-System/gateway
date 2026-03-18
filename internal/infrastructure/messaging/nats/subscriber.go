package nats

import (
	"fmt"

	"github.com/fercho/school-tracking/services/gateway/internal/core/ports/resources"
	"github.com/nats-io/nats.go"
)

const (
	durableConsumer = "gateway-fleet"
)

type subscriber struct {
	js            nats.JetStreamContext
	subscriptions []*nats.Subscription
}

// NewSubscriber creates a JetStream subscriber bound to the given NATS connection.
func NewSubscriber(nc *nats.Conn) (*subscriber, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}
	return &subscriber{js: js}, nil
}

// Subscribe creates a durable JetStream push subscription for the given subject.
func (s *subscriber) Subscribe(subject string, handler resources.EventHandler) error {
	sub, err := s.js.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Subject, msg.Data)
		_ = msg.Ack()
	}, nats.Durable(durableConsumer), nats.DeliverAll())
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}
	s.subscriptions = append(s.subscriptions, sub)
	return nil
}

// Close drains all active subscriptions.
func (s *subscriber) Close() {
	for _, sub := range s.subscriptions {
		_ = sub.Drain()
	}
}
