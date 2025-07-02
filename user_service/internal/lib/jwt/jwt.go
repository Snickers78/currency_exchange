package jwt

import (
	"fmt"
	"time"
	"user_service/internal/domain/model"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user *model.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("cannot create token: %w", err)
	}

	return signedToken, nil

}
