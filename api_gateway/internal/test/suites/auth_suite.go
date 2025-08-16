package test

import (
	"api_gateway/internal/auth"
	"api_gateway/internal/config"
	"context"
	"testing"
	"time"
)

type AuthSuite struct {
	Test   *testing.T
	Client *auth.AuthClient
	Cfg    *config.Config
}

func NewAuthSuite(t *testing.T) (context.Context, *AuthSuite) {
	t.Helper()
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cfg := config.LoadConfig(Test)

	client := auth.NewAuthClient(cfg)

	return ctx, &AuthSuite{
		Test:   t,
		Client: client,
		Cfg:    cfg,
	}
}
