package auth

import (
	"context"
	"errors"
	"time"

	oauthloginv1 "github.com/mandacode-com/accounts-proto/auth/login/oauth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewOAuthLoginClient(addr string) (oauthloginv1.OAuthLoginServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	if conn == nil {
		return nil, nil, errors.New("gRPC client connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	if healthClient == nil {
		return nil, nil, errors.New("gRPC health client is nil")
	}
	healthResp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
	if err != nil {
		return nil, nil, err
	}
	if healthResp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return nil, nil, errors.New("token service is not serving")
	}

	client := oauthloginv1.NewOAuthLoginServiceClient(conn)
	if client == nil {
		return nil, nil, errors.New("gRPC client is nil")
	}

	return client, conn, nil
}
