package env

import (
	"github.com/caarlos0/env/v10"
	"go.uber.org/fx"
)

type Config struct {
	Port            string `env:"PORT" envDefault:"8000"`
	Environment     string `env:"ENVIRONMENT" envDefault:"development"`
	AuthServiceURL  string `env:"AUTH_SERVICE_URL" envDefault:"http://localhost:8080"`
	AuthGRPCURL     string `env:"AUTH_GRPC_URL" envDefault:"localhost:9090"`
	FleetServiceURL        string `env:"FLEET_SERVICE_URL" envDefault:"localhost:9090"`
	TripServiceURL         string `env:"TRIP_SERVICE_URL" envDefault:"localhost:9092"`
	NotificationServiceURL string `env:"NOTIFICATION_SERVICE_URL" envDefault:"localhost:9095"`
	JWTSecret              string `env:"JWT_SECRET" envDefault:"dev-secret-change-in-prod"`
	NatsURL         string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

var Module = fx.Provide(NewConfig)
