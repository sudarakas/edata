package user

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sudarakas/edata/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, user types.User) error {
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (firstName, lastName, email, password) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (types.User, error) {
	var user types.User
	stmt, err := s.db.PrepareContext(ctx, "SELECT id, firstName, lastName, email, password, createdAt, updatedAt FROM users WHERE email = $1")
	if err != nil {
		return types.User{}, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, nil
		}
		return types.User{}, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}
	return user, nil
}

func (s *Store) GetUserByID(ctx context.Context, id string) (types.User, error) {
	var user types.User
	stmt, err := s.db.PrepareContext(ctx, "SELECT id, firstName, lastName, email, password, createdAt, updatedAt FROM users WHERE id = $1")
	if err != nil {
		return types.User{}, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, nil
		}
		return types.User{}, fmt.Errorf("failed to get user by id %s: %w", id, err)
	}
	return user, nil
}
