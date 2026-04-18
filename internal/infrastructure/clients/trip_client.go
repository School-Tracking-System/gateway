package clients

import (
	"context"
	"fmt"

	pb "github.com/fercho/school-tracking/proto/gen/trip/v1"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TripClients groups all gRPC service clients for the Trip service.
type TripClients struct {
	fx.Out

	TripService pb.TripServiceClient
}

// NewTripClient opens a single gRPC connection to the Trip service.
func NewTripClient(lc fx.Lifecycle, cfg *env.Config, log *zap.Logger) (TripClients, error) {
	addr := cfg.TripServiceURL

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return TripClients{}, fmt.Errorf("failed to connect to trip service at %s: %w", addr, err)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("Connected to Trip gRPC service", zap.String("addr", addr))
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("Closing Trip gRPC connection")
			return conn.Close()
		},
	})

	return TripClients{
		TripService: pb.NewTripServiceClient(conn),
	}, nil
}
