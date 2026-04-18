package main

import (
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/api/handlers"
	"github.com/fercho/school-tracking/services/gateway/internal/infrastructure/clients"
	messaging "github.com/fercho/school-tracking/services/gateway/internal/infrastructure/messaging/nats"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	"github.com/fercho/school-tracking/services/gateway/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func AppModule() fx.Option {
	return fx.Options(
		env.Module,
		logger.Module,
		messaging.Module,
		clients.Module, // Provides FleetClient (grpc)
		fx.Provide(handlers.NewFleetHandler),
		fx.Provide(handlers.NewTripHandler),		fx.Provide(handlers.NewNotificationHandler),
		fx.Provide(api.NewRouter),
		fx.Invoke(api.StartHTTPServer),
		fx.Invoke(func(l *zap.Logger, cfg *env.Config) {
			l.Info("Gateway service initialized", zap.String("port", cfg.Port), zap.String("env", cfg.Environment))
		}),
	)
}
