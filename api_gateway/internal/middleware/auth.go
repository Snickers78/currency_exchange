package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	ContextKeyEmail  ContextKey = "email"
	ContextKeyUserID ContextKey = "id"
)

type JWTClaims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func writeUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
}

func IsAuthed(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			writeUnauthorized(c)
			return
		}

		updatedToken := strings.TrimPrefix(token, "Bearer ")

		parsedToken, err := jwt.ParseWithClaims(updatedToken, &JWTClaims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil {
			writeUnauthorized(c)
			return
		}

		// Проверяем валидность токена
		if !parsedToken.Valid {
			writeUnauthorized(c)
			return
		}

		// Извлекаем claims
		claims, ok := parsedToken.Claims.(*JWTClaims)
		if !ok {
			writeUnauthorized(c)
			return
		}

		// Добавляем данные пользователя в контекст
		ctx := context.WithValue(c.Request.Context(), ContextKeyEmail, claims.Email)
		ctx = context.WithValue(ctx, ContextKeyUserID, claims.ID)

		// Обновляем контекст запроса
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
