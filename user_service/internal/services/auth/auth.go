package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"
	"user_service/infra/metrics"
	"user_service/internal/config"
	errs "user_service/internal/domain/errors"
	"user_service/internal/domain/model"
	"user_service/internal/lib/jwt"
	"user_service/internal/lib/logger"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	config   *config.Config
	storage  Storage
	tokenTTL time.Duration
	logger   *slog.Logger
}

type Storage interface {
	CreateUser(ctx context.Context, email string, password []byte) (int64, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
}

// New returns new instance of Auth
func NewAuthService(storage Storage, tokenTTL time.Duration, config *config.Config, logger *slog.Logger) *Auth {
	return &Auth{
		storage:  storage,
		tokenTTL: tokenTTL,
		config:   config,
		logger:   logger,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	metrics.RequestCount.Inc()
	start := time.Now()

	user, err := a.storage.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, errs.UserNotFound) {
			logEntry := logger.NewAuthLog("error", "login_wrong_credentials", logger.WithAuthEmail(email), logger.WithAuthError(errs.WrongCredentials.Error()))
			if msg, err := json.Marshal(logEntry); err == nil {
				a.logger.Error("error", "auth_service", string(msg))
			}
			metrics.ErrorCount.Inc()
			return "", errs.WrongCredentials
		}
		logEntry := logger.NewAuthLog("error", "login_get_user_failed", logger.WithAuthEmail(email), logger.WithAuthError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			a.logger.Error("error", "auth_service", string(msg))
		}
		return "", errs.BadRequest
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		logEntry := logger.NewAuthLog("error", "login_wrong_credentials", logger.WithAuthEmail(email), logger.WithAuthError(errs.WrongCredentials.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			a.logger.Error("error", "auth_service", string(msg))
		}
		metrics.ErrorCount.Inc()
		return "", errs.WrongCredentials
	}

	token, err := jwt.NewToken(user, a.tokenTTL, a.config.Secret)
	if err != nil {
		logEntry := logger.NewAuthLog("error", "login_token_failed", logger.WithAuthEmail(email), logger.WithAuthError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			a.logger.Error("error", "auth_service", string(msg))
		}
		metrics.ErrorCount.Inc()
		return "", errs.InternalError
	}

	metrics.ResponseTimeSeconds.Observe(time.Since(start).Seconds())
	return token, nil
}

// Register new user. If user already exists returns error
func (a *Auth) Register(ctx context.Context, email, password string) (int64, error) {
	metrics.RequestCount.Inc()
	start := time.Now()

	user, err := a.storage.GetUser(ctx, email)
	if err != nil {
		if !errors.Is(err, errs.UserNotFound) {
			logEntry := logger.NewAuthLog("error", "register_check_user_failed", logger.WithAuthEmail(email), logger.WithAuthError(err.Error()))
			if msg, err := json.Marshal(logEntry); err == nil {
				a.logger.Error("error", "auth_service", string(msg))
			}
			metrics.ErrorCount.Inc()
			return 0, err
		}
	} else if user != nil {
		logEntry := logger.NewAuthLog("error", "register_user_exists", logger.WithAuthEmail(email), logger.WithAuthError(errs.UserAlreadyExists.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			a.logger.Error("error", "auth_service", string(msg))
		}
		return 0, errs.UserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logEntry := logger.NewAuthLog("error", "register_hash_failed", logger.WithAuthEmail(email), logger.WithAuthError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			a.logger.Error("error", "auth_service", string(msg))
		}
		metrics.ErrorCount.Inc()
		return 0, err
	}

	uid, err := a.storage.CreateUser(ctx, email, hash)
	if err != nil {
		logEntry := logger.NewAuthLog("error", "register_create_user_failed", logger.WithAuthEmail(email), logger.WithAuthError(err.Error()))
		if msg, err := json.Marshal(logEntry); err == nil {
			a.logger.Error("error", "auth_service", string(msg))
		}
		metrics.ErrorCount.Inc()
		return 0, err
	}

	metrics.ResponseTimeSeconds.Observe(time.Since(start).Seconds())

	return uid, nil
}
