package auth

import (
	"api_gateway/internal/config"
	ssov1 "api_gateway/internal/gen/auth/proto"
	"context"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client ssov1.AuthClient
}

func NewAuthCLient(cfg *config.Config) *AuthClient {
	conn, err := grpc.NewClient(net.JoinHostPort("localhost", strconv.Itoa(cfg.AuthServicePort)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &AuthClient{
		client: ssov1.NewAuthClient(conn),
	}
}

func (c *AuthClient) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	return c.client.Login(ctx, req)
}

func (c *AuthClient) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	return c.client.Register(ctx, req)
}
