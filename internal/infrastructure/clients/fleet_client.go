package clients

import (
	"context"
	"fmt"

	pb "github.com/fercho/school-tracking/proto/gen/fleet/v1"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FleetClient wraps the gRPC client for the Fleet service
type FleetClient struct {
	Vehicles pb.VehicleServiceClient
	conn     *grpc.Server
}

func NewFleetClient(lc fx.Lifecycle, cfg *env.Config, log *zap.Logger) (pb.VehicleServiceClient, error) {
	addr := cfg.FleetServiceURL
	
	// Create connection with insecure credentials (for now)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to fleet service: %w", err)
	}

	client := pb.NewVehicleServiceClient(conn)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Connecting to Fleet service", zap.String("addr", addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing connection to Fleet service")
			return conn.Close()
		},
	})

	return client, nil
}

var Module = fx.Provide(NewFleetClient)
