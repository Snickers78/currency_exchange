package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	errs "user_service/internal/domain/errors"
	"user_service/internal/domain/model"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) *Storage {

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		panic(err)
	}

	return &Storage{db: db}
}

func (s *Storage) CreateUser(ctx context.Context, email string, password []byte) (int64, error) {
	var userId int64
	err := s.db.QueryRowContext(ctx, "INSERT INTO users (email, hashedpassword) VALUES ($1, $2) RETURNING id", email, password).Scan(&userId)
	if err != nil {
		log.Printf("Failed to create user with email %s: %v", email, err)
		return 0, err
	}

	return userId, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (*model.User, error) {
	stmt, err := s.db.PrepareContext(ctx, "SELECT id, email, hashedpassword FROM users WHERE email = $1")
	if err != nil {
		log.Printf("Failed to prepare statement for getting user: %v", err)
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	if err := stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.HashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.UserNotFound
		}
		return nil, err
	}

	return &user, nil
}
