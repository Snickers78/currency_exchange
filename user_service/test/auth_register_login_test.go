package tests

import (
	"testing"
	"time"
	ssov1 "user_service/gen"
	"user_service/test/suite"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLogin(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 10)

	req, err := st.Client.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, req.GetUserId())

	resp, err := st.Client.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
	})

	loginTime := time.Now()

	require.NoError(t, err)
	require.NotEmpty(t, resp.GetToken())

	token, err := jwt.Parse(resp.GetToken(), func(token *jwt.Token) (interface{}, error) {
		return []byte(st.Cfg.Secret), nil
	})
	require.NoError(t, err)
	require.NotNil(t, token)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, req.GetUserId(), int64(claims["id"].(float64)))

	assert.InDelta(t, loginTime.Add(st.Cfg.TockenTTL).Unix(), int64(claims["exp"].(float64)), 1)
}
