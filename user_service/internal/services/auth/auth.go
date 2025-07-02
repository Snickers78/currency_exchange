package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"user_service/internal/config"
	errs "user_service/internal/domain/errors"
	"user_service/internal/domain/model"
	"user_service/internal/lib/jwt"
	"user_service/internal/metrics"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	config   *config.Config
	log      *slog.Logger
	storage  Storage
	tokenTTL time.Duration
}

type Storage interface {
	CreateUser(ctx context.Context, email string, password []byte) (int64, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
}

// New returns new instance of Auth
func NewAuthService(log *slog.Logger, storage Storage, tokenTTL time.Duration, config *config.Config) *Auth {
	return &Auth{log: log,
		storage:  storage,
		tokenTTL: tokenTTL,
		config:   config,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	metrics.RequestCount.Inc()
	start := time.Now()
	user, err := a.storage.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, errs.UserNotFound) {
			a.log.Warn("Cannot login user: wrong credentials")
			metrics.ErrorCount.Inc()
			return "", errs.WrongCredentials
		}
		a.log.Warn("Cannot get user: ", err)
		return "", errs.BadRequest
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		a.log.Warn("Cannot login user: wrong credentials")
		metrics.ErrorCount.Inc()
		return "", errs.WrongCredentials
	}

	token, err := jwt.NewToken(user, a.tokenTTL, a.config.Secret)
	if err != nil {
		a.log.Error("Failed to create token: ", err)
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
			a.log.Error("Failed to check if user exists: ", err)
			metrics.ErrorCount.Inc()
			return 0, err
		}
	} else if user != nil {
		a.log.Warn("Cannot register new user: user already exists", "email", email)
		return 0, errs.UserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("Failed to hash password: ", err)
		metrics.ErrorCount.Inc()
		return 0, err
	}

	uid, err := a.storage.CreateUser(ctx, email, hash)
	if err != nil {
		a.log.Error("Failed to create user: ", err)
		metrics.ErrorCount.Inc()
		return 0, err
	}

	metrics.ResponseTimeSeconds.Observe(time.Since(start).Seconds())
	return uid, nil
}
