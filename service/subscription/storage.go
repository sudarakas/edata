package subscription

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

// Common errors
var (
	ErrNotFound      = errors.New("subscription not found")
	ErrDuplicateCode = errors.New("subscription code already exists")
)

type Store struct {
	db *sql.DB
	// Pre-prepared statements
	createSubscriptionStmt          *sql.Stmt
	getSubscriptionByCodeStmt       *sql.Stmt
	getSubscriptionByPriceRangeStmt *sql.Stmt
	updateSubscriptionStmt          *sql.Stmt
	deleteSubscriptionStmt          *sql.Stmt
}

// Queries
const (
	createSubscriptionQuery = `
		INSERT INTO subscriptions (code, name, monthly_price, active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, code, name, monthly_price, active, created_at, updated_at
	`

	getSubscriptionByCodeQuery = `
		SELECT id, code, name, monthly_price, active, created_at, updated_at
		FROM subscriptions
		WHERE code = $1
	`

	getSubscriptionByPriceRangeQuery = `
		SELECT id, code, name, monthly_price, active, created_at, updated_at
		FROM subscriptions
		WHERE monthly_price >= $1 AND monthly_price <= $2
		ORDER BY monthly_price
	`

	updateSubscriptionQuery = `
		UPDATE subscriptions
		SET code = $1, name = $2, monthly_price = $3, active = $4
		WHERE id = $5
		RETURNING id, code, name, monthly_price, active, created_at, updated_at
	`

	deleteSubscriptionQuery = `
		DELETE FROM subscriptions
		WHERE id = $1
	`
)

// Initializes the Store with prepared statements
func NewStore(db *sql.DB) (*Store, error) {
	store := &Store{
		db: db,
	}

	if err := store.prepareStatements(); err != nil {
		return nil, fmt.Errorf("failed to prepare statements: %w", err)
	}

	return store, nil
}

// Prepare statements
func (s *Store) prepareStatements() error {
	var err error

	if s.createSubscriptionStmt, err = s.db.Prepare(createSubscriptionQuery); err != nil {
		return fmt.Errorf("failed to prepare createSubscription statement: %w", err)
	}

	if s.getSubscriptionByCodeStmt, err = s.db.Prepare(getSubscriptionByCodeQuery); err != nil {
		return fmt.Errorf("failed to prepare getSubscriptionByCode statement: %w", err)
	}

	if s.getSubscriptionByPriceRangeStmt, err = s.db.Prepare(getSubscriptionByPriceRangeQuery); err != nil {
		return fmt.Errorf("failed to prepare getSubscriptionByPriceRange statement: %w", err)
	}

	if s.updateSubscriptionStmt, err = s.db.Prepare(updateSubscriptionQuery); err != nil {
		return fmt.Errorf("failed to prepare updateSubscription statement: %w", err)
	}

	if s.deleteSubscriptionStmt, err = s.db.Prepare(deleteSubscriptionQuery); err != nil {
		return fmt.Errorf("failed to prepare deleteSubscription statement: %w", err)
	}

	return nil
}

// Close the prepared statements
func (s *Store) Close() error {
	var errs []error

	if err := s.createSubscriptionStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close createSubscription statement: %w", err))
	}
	if err := s.getSubscriptionByCodeStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close getSubscriptionByCode statement: %w", err))
	}
	if err := s.getSubscriptionByPriceRangeStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close getSubscriptionByPriceRange statement: %w", err))
	}
	if err := s.updateSubscriptionStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close updateSubscription statement: %w", err))
	}
	if err := s.deleteSubscriptionStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close deleteSubscription statement: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close statements: %v", errs)
	}

	return nil
}
