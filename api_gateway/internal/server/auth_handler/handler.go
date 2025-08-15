package handler

import (
	"api_gateway/infra/kafka"
	"api_gateway/internal/auth"
	ssov1 "api_gateway/internal/gen/auth/proto"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	authClient *auth.AuthClient
	hook       *kafka.KafkaHook
}

func NewAuthHandler(router *gin.Engine, authClient *auth.AuthClient, hook *kafka.KafkaHook, secret string) {
	handler := &Auth{
		authClient: authClient,
		hook:       hook,
	}

	router.POST("/", handler.Login)
	router.POST("/register", handler.Register)
}

func (a *Auth) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := NewAuthLog("error", "login_invalid_request", WithAuthError(err.Error()))
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
		logEntry := NewAuthLog("error", "login_failed", WithAuthEmail(req.Email), WithAuthError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			a.hook.Fire(string(msg))
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
		return
	}

	logEntry := NewAuthLog("info", "login_success", WithAuthEmail(req.Email))
	if msg, err := json.Marshal(logEntry); err == nil {
		a.hook.Fire(string(msg))
	}
	c.JSON(http.StatusOK, gin.H{"token": resp.Token})
}

func (a *Auth) Register(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := NewAuthLog("error", "register_invalid_request", WithAuthError(err.Error()))
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
		logEntry := NewAuthLog("error", "register_failed", WithAuthEmail(req.Email), WithAuthError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			a.hook.Fire(string(msg))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed"})
		return
	}

	logEntry := NewAuthLog("info", "register_success", WithAuthEmail(req.Email), WithAuthUserID(resp.GetUserId()))
	if msg, err := json.Marshal(logEntry); err == nil {
		if err := a.hook.Fire(string(msg)); err != nil {
			errLog := NewAuthLog("error", "register_kafka_error", WithAuthUserID(resp.GetUserId()), WithAuthError(err.Error()))
			if emsg, e := json.Marshal(errLog); e == nil {
				a.hook.Fire(string(emsg))
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId})
}
