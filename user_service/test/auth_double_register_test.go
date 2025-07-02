package tests

import (
	"testing"
	ssov1 "user_service/gen"
	"user_service/test/suite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDoubleRegister(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email := "register@mail.ru"
	password := "pass333"

	req, err := st.Client.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, req.GetUserId())

	req, err = st.Client.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	assert.Equal(t, status.Error(codes.AlreadyExists, "user already exists"), err)
}
