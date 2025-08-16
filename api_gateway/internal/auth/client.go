package auth

import (
	"api_gateway/internal/config"
	ssov1 "api_gateway/internal/gen/auth/proto"
	"context"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client ssov1.AuthClient
}

func NewAuthClient(cfg *config.Config) *AuthClient {
	// systemRoots, err := x509.SystemCertPool()
	// if err != nil {
	// 	return nil
	// }

	// tlsCreds := credentials.NewTLS(&tls.Config{
	// 	RootCAs: systemRoots,
	// })

	conn, err := grpc.NewClient(net.JoinHostPort("localhost", strconv.Itoa(cfg.AuthServicePort)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &AuthClient{
		client: ssov1.NewAuthClient(conn),
	}
}

func (c *AuthClient) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	resp, err := c.client.Login(ctx, req)
	if err != nil {
		log.Println(err)
	}
	return resp, nil
}

func (c *AuthClient) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	resp, err := c.client.Register(ctx, req)
	if err != nil {
		log.Println(err)
	}
	return resp, nil
}
