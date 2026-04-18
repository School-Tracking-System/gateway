package clients

import (
	"context"
	"fmt"

	pb "github.com/fercho/school-tracking/proto/gen/notification/v1"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NotificationClients groups all gRPC service clients for the Notification service.
type NotificationClients struct {
	fx.Out

	NotificationService pb.NotificationServiceClient
}

// NewNotificationClient opens a single gRPC connection to the Notification service.
func NewNotificationClient(lc fx.Lifecycle, cfg *env.Config, log *zap.Logger) (NotificationClients, error) {
	addr := cfg.NotificationServiceURL

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return NotificationClients{}, fmt.Errorf("failed to connect to notification service at %s: %w", addr, err)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("Connected to Notification gRPC service", zap.String("addr", addr))
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("Closing Notification gRPC connection")
			return conn.Close()
		},
	})

	return NotificationClients{
		NotificationService: pb.NewNotificationServiceClient(conn),
	}, nil
}
