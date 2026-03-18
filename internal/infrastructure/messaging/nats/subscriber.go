package nats

import (
	"fmt"
	"strings"

	"github.com/fercho/school-tracking/services/gateway/internal/core/ports/resources"
	"github.com/nats-io/nats.go"
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
// The durable consumer name is derived from the subject so each subscription gets
// its own unique consumer (e.g. "fleet.vehicle.created" → "gateway-fleet-vehicle-created").
func (s *subscriber) Subscribe(subject string, handler resources.EventHandler) error {
	durable := "gateway-" + strings.ReplaceAll(subject, ".", "-")
	sub, err := s.js.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Subject, msg.Data)
		_ = msg.Ack()
	}, nats.Durable(durable), nats.DeliverAll())
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
