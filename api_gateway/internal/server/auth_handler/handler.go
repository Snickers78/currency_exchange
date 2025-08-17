package handler

import (
	"api_gateway/internal/auth"
	ssov1 "api_gateway/internal/gen/auth/proto"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	authClient *auth.AuthClient
	logger     *slog.Logger
}

func NewAuthHandler(router *gin.Engine, authClient *auth.AuthClient, logger *slog.Logger, secret string) {
	handler := &Auth{
		authClient: authClient,
		logger:     logger,
	}

	router.POST("/login", handler.Login)
	router.POST("/register", handler.Register)
}

func (a *Auth) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := NewAuthLog("error", "login_invalid_request", WithAuthError(err.Error()))
		a.logger.Error("Error", "auth", logEntry)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := a.authClient.Login(c.Request.Context(), &ssov1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logEntry := NewAuthLog("error", "login_failed", WithAuthEmail(req.Email), WithAuthError(err.Error()))
		a.logger.Error("Error", "auth", logEntry)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed"})
		return
	}

	logEntry := NewAuthLog("info", "login_success", WithAuthEmail(req.Email))
	a.logger.Info("Success", "auth", logEntry)
	c.JSON(http.StatusOK, gin.H{"token": resp.Token})
}

func (a *Auth) Register(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry := NewAuthLog("error", "register_invalid_request", WithAuthError(err.Error()))
		a.logger.Error("Error", "auth", logEntry)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := a.authClient.Register(c.Request.Context(), &ssov1.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logEntry := NewAuthLog("error", "register_failed", WithAuthEmail(req.Email), WithAuthError(err.Error()))
		a.logger.Error("Error", "auth", logEntry)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed"})
		return
	}

	logEntry := NewAuthLog("info", "register_success", WithAuthEmail(req.Email), WithAuthUserID(resp.GetUserId()))
	a.logger.Info("Success", "auth", logEntry)
	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId})
}
