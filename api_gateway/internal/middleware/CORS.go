package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	allowedMethods := "GET, POST, OPTIONS"
	allowedHeaders := "Content-Type, Authorization, X-Requested-With"
	allowOrigin := "*"

	return func(c *gin.Context) {
		// Минимальная логика
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", allowedMethods)
		c.Header("Access-Control-Allow-Headers", allowedHeaders)

		// Для OPTIONS запросов - сразу отвечаем
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
