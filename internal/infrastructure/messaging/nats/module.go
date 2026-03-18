package nats

import (
	"context"
	"fmt"

	"github.com/fercho/school-tracking/services/gateway/internal/core/ports/resources"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/messaging/nats/handlers"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	subjectVehicleCreated = "fleet.vehicle.created"
	subjectVehicleUpdated = "fleet.vehicle.updated"
)

// NewConnection establishes a NATS connection using the URL from config.
func NewConnection(lc fx.Lifecycle, cfg *env.Config, log *zap.Logger) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.NatsURL,
		nats.Name("gateway-service"),
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS at %s: %w", cfg.NatsURL, err)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("Connected to NATS", zap.String("url", cfg.NatsURL))
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("Closing NATS connection")
			nc.Drain()
			return nil
		},
	})

	return nc, nil
}

// NewEventSubscriber wraps the NATS subscriber as the EventSubscriber port.
func NewEventSubscriber(nc *nats.Conn) (resources.EventSubscriber, error) {
	return NewSubscriber(nc)
}

// RegisterFleetSubscriptions wires fleet event subjects to their handlers.
func RegisterFleetSubscriptions(lc fx.Lifecycle, sub resources.EventSubscriber, handler *handlers.FleetEventHandler, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := sub.Subscribe(subjectVehicleCreated, handler.HandleVehicleCreated); err != nil {
				return fmt.Errorf("failed to subscribe to %s: %w", subjectVehicleCreated, err)
			}
			if err := sub.Subscribe(subjectVehicleUpdated, handler.HandleVehicleUpdated); err != nil {
				return fmt.Errorf("failed to subscribe to %s: %w", subjectVehicleUpdated, err)
			}
			log.Info("Subscribed to fleet events",
				zap.String("subjects", subjectVehicleCreated+", "+subjectVehicleUpdated),
			)
			return nil
		},
		OnStop: func(_ context.Context) error {
			sub.Close()
			return nil
		},
	})
}

// Module provides NATS connection, subscriber, and fleet event handlers to the fx dependency graph.
var Module = fx.Module("messaging.nats",
	fx.Provide(NewConnection),
	fx.Provide(NewEventSubscriber),
	fx.Provide(handlers.NewFleetEventHandler),
	fx.Invoke(RegisterFleetSubscriptions),
)
