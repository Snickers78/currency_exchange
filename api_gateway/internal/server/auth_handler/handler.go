package handler

import (
	"api_gateway/internal/auth"
	ssov1 "api_gateway/internal/gen/auth/proto"
	"api_gateway/internal/kafka"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	authClient *auth.AuthClient
	hook       *kafka.KafkaHook
}

func NewAuthHandler(router *gin.Engine, authClient *auth.AuthClient, hook *kafka.KafkaHook) *Auth {
	handler := &Auth{
		authClient: authClient,
		hook:       hook,
	}
	router.POST("/login", handler.Login)
	router.POST("/register", handler.Register)
	return handler
}

func (a *Auth) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := AuthLog{
			Level: "error",
			Event: "login_invalid_request",
			Error: err.Error(),
			Time:  time.Now().Format(time.RFC3339),
		}
		if msg, err := json.Marshal(logEntry); err == nil {
			a.hook.Fire(string(msg))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := a.authClient.Login(c.Request.Context(), &ssov1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logEntry := AuthLog{
			Level: "error",
			Event: "login_failed",
			Email: req.Email,
			Error: err.Error(),
			Time:  time.Now().Format(time.RFC3339),
		}
		if msg, err := json.Marshal(logEntry); err == nil {
			a.hook.Fire(string(msg))
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
		return
	}

	logEntry := AuthLog{
		Level: "info",
		Event: "login_success",
		Email: req.Email,
		Time:  time.Now().Format(time.RFC3339),
	}
	if msg, err := json.Marshal(logEntry); err == nil {
		a.hook.Fire(string(msg))
	}
	c.JSON(http.StatusOK, gin.H{"token": resp.Token})
}

func (a *Auth) Register(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := AuthLog{
			Level: "error",
			Event: "register_invalid_request",
			Error: err.Error(),
			Time:  time.Now().Format(time.RFC3339),
		}
		if msg, err := json.Marshal(logEntry); err == nil {
			a.hook.Fire(string(msg))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := a.authClient.Register(c.Request.Context(), &ssov1.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logEntry := AuthLog{
			Level: "error",
			Event: "register_failed",
			Email: req.Email,
			Error: err.Error(),
			Time:  time.Now().Format(time.RFC3339),
		}
		if msg, err := json.Marshal(logEntry); err == nil {
			a.hook.Fire(string(msg))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed"})
		return
	}

	logEntry := AuthLog{
		Level:  "info",
		Event:  "register_success",
		Email:  req.Email,
		UserID: resp.GetUserId(),
		Time:   time.Now().Format(time.RFC3339),
	}
	if msg, err := json.Marshal(logEntry); err == nil {
		if err := a.hook.Fire(string(msg)); err != nil {
			errLog := AuthLog{
				Level:  "error",
				Event:  "register_kafka_error",
				UserID: resp.GetUserId(),
				Error:  err.Error(),
				Time:   time.Now().Format(time.RFC3339),
			}
			if emsg, e := json.Marshal(errLog); e == nil {
				a.hook.Fire(string(emsg))
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId})
}
