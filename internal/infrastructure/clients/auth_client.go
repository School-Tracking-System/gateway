package clients

import (
	"context"
	"fmt"

	authpb "github.com/fercho/school-tracking/proto/gen/auth/v1"
	"github.com/fercho/school-tracking/services/gateway/pkg/env"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(lc fx.Lifecycle, cfg *env.Config, log *zap.Logger) (authpb.AuthServiceClient, error) {
	addr := cfg.AuthGRPCURL

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth gRPC service: %w", err)
	}

	client := authpb.NewAuthServiceClient(conn)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Connecting to Auth gRPC service", zap.String("addr", addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing connection to Auth gRPC service")
			return conn.Close()
		},
	})

	return client, nil
}
