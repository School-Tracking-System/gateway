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

// FleetClients groups all gRPC service clients for the Fleet service.
// Using fx.Out allows fx to inject each client type independently.
type FleetClients struct {
	fx.Out

	Vehicles  pb.VehicleServiceClient
	Schools   pb.SchoolServiceClient
	Drivers   pb.DriverServiceClient
	Students  pb.StudentServiceClient
	Guardians pb.GuardianServiceClient
	Routes    pb.RouteServiceClient
}

// NewFleetClient opens a single gRPC connection to the Fleet service and
// creates clients for all its sub-services.
func NewFleetClient(lc fx.Lifecycle, cfg *env.Config, log *zap.Logger) (FleetClients, error) {
	addr := cfg.FleetServiceURL

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return FleetClients{}, fmt.Errorf("failed to connect to fleet service at %s: %w", addr, err)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("Connected to Fleet gRPC service", zap.String("addr", addr))
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("Closing Fleet gRPC connection")
			return conn.Close()
		},
	})

	return FleetClients{
		Vehicles:  pb.NewVehicleServiceClient(conn),
		Schools:   pb.NewSchoolServiceClient(conn),
		Drivers:   pb.NewDriverServiceClient(conn),
		Students:  pb.NewStudentServiceClient(conn),
		Guardians: pb.NewGuardianServiceClient(conn),
		Routes:    pb.NewRouteServiceClient(conn),
	}, nil
}
