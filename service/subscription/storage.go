package subscription

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/sudarakas/edata/types"
)

// Common errors
var (
	ErrNotFound      = errors.New("subscription not found")
	ErrDuplicateCode = errors.New("subscription code already exists")
)

type Store struct {
	db *sql.DB
	// Pre-prepared statements
	createSubscriptionStmt    *sql.Stmt
	getSubscriptionById       *sql.Stmt
	getSubscriptionByCodeStmt *sql.Stmt
	GetAllSubscriptionsStmt   *sql.Stmt
	updateSubscriptionStmt    *sql.Stmt
	deleteSubscriptionStmt    *sql.Stmt
	filterSubscriptions       *sql.Stmt
}

// Queries
const (
	createSubscriptionQuery = `
		INSERT INTO subscriptions (code, name, description, price, data_limit, validity_period)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, code, name, description, price, data_limit, validity_period AS validity, created_at, updated_at
	`

	getSubscriptionByIdQuery = `
		SELECT id, code, name, description, price, data_limit, validity_period AS validity, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	getSubscriptionByCodeQuery = `
		SELECT id, code, name, description, price, data_limit, validity_period AS validity, created_at, updated_at
		FROM subscriptions
		WHERE code = $1
	`

	getAllSubscriptionsQuery = `
		SELECT id, code, name, description, price, data_limit, validity_period AS validity, created_at, updated_at
		FROM subscriptions
		WHERE deleted_at IS NULL
		ORDER BY 
			CASE $1 
				WHEN 'id' THEN id::text
				WHEN 'code' THEN code
				WHEN 'name' THEN name
				WHEN 'created_at' THEN created_at::text
				ELSE id::text 
			END ASC
		LIMIT $2 OFFSET $3
	`

	updateSubscriptionQuery = `
		UPDATE subscriptions
		SET code = $1, name = $2, description = $3, price = $4, data_limit = $5, validity_period = $6
		WHERE id = $7
		RETURNING id, code, name, description, price, data_limit, validity_period AS validity, created_at, updated_at
	`

	deleteSubscriptionQuery = `
		UPDATE subscriptions 
		SET 
			deleted_at = NOW() 
		WHERE id = $1
	`

	filterSubscriptionsQuery = `
		SELECT id, code, name, description, price, data_limit, validity_period AS validity, created_at, updated_at
		FROM subscriptions
		WHERE
			($1::TEXT IS NULL OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%') AND
			($2::FLOAT IS NULL OR price >= $2) AND
			($3::FLOAT IS NULL OR price <= $3) AND
			($4::INT IS NULL OR data_limit >= $4) AND
			($5::INT IS NULL OR data_limit <= $5) AND
			($6::INT IS NULL OR validity_period >= $6) AND
			($7::INT IS NULL OR validity_period <= $7)
		LIMIT $8 OFFSET $9
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
	if s.getSubscriptionById, err = s.db.Prepare(getSubscriptionByIdQuery); err != nil {
		return fmt.Errorf("failed to prepare getSubscriptionById statement: %w", err)
	}
	if s.getSubscriptionByCodeStmt, err = s.db.Prepare(getSubscriptionByCodeQuery); err != nil {
		return fmt.Errorf("failed to prepare getSubscriptionByCode statement: %w", err)
	}
	if s.GetAllSubscriptionsStmt, err = s.db.Prepare(getAllSubscriptionsQuery); err != nil {
		return fmt.Errorf("failed to prepare getAllSubscriptions statement: %w", err)
	}
	if s.updateSubscriptionStmt, err = s.db.Prepare(updateSubscriptionQuery); err != nil {
		return fmt.Errorf("failed to prepare updateSubscription statement: %w", err)
	}
	if s.deleteSubscriptionStmt, err = s.db.Prepare(deleteSubscriptionQuery); err != nil {
		return fmt.Errorf("failed to prepare deleteSubscription statement: %w", err)
	}
	if s.filterSubscriptions, err = s.db.Prepare(filterSubscriptionsQuery); err != nil {
		return fmt.Errorf("failed to prepare filterSubscriptions statement: %w", err)
	}

	return nil
}

// Close the prepared statements
func (s *Store) Close() error {
	var errs []error

	if err := s.createSubscriptionStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close createSubscription statement: %w", err))
	}
	if err := s.getSubscriptionById.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close getSubscriptionById statement: %w", err))
	}
	if err := s.getSubscriptionByCodeStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close getSubscriptionByCode statement: %w", err))
	}
	if err := s.GetAllSubscriptionsStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close getAllSubscriptions statement: %w", err))
	}
	if err := s.updateSubscriptionStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close updateSubscription statement: %w", err))
	}
	if err := s.deleteSubscriptionStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close deleteSubscription statement: %w", err))
	}
	if err := s.filterSubscriptions.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close filterSubscriptions statement: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close statements: %v", errs)
	}

	return nil
}

// CreateSubscription creates a new subscription
func (s *Store) CreateSubscription(ctx context.Context, subscription types.Subscription) (types.Subscription, error) {
	err := s.createSubscriptionStmt.QueryRow(
		subscription.Code,
		subscription.Plan,
		subscription.Description,
		subscription.Price,
		subscription.DataLimit,
		subscription.Validity,
	).Scan(
		&subscription.SubscriptionID,
		&subscription.Code,
		&subscription.Plan,
		&subscription.Description,
		&subscription.Price,
		&subscription.DataLimit,
		&subscription.Validity,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		if isPgDuplicateKeyError(err) {
			return types.Subscription{}, ErrDuplicateCode
		}
		return types.Subscription{}, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

// isPgDuplicateKeyError checks if the error is a PostgreSQL unique violation
func isPgDuplicateKeyError(err error) bool {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return pqErr.Code == "23505" // unique_violation error code
}

// GetSubscriptionByID retrieves a subscription by its ID
func (s *Store) GetSubscriptionByID(ctx context.Context, id int) (types.Subscription, error) {
	var subscription types.Subscription
	var deletedAt sql.NullTime

	err := s.getSubscriptionById.QueryRow(id).Scan(
		&subscription.SubscriptionID,
		&subscription.Code,
		&subscription.Plan,
		&subscription.Description,
		&subscription.Price,
		&subscription.DataLimit,
		&subscription.Validity,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&deletedAt,
	)

	// Check if the subscription is not deleted
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Subscription{}, ErrNotFound
		}
		return types.Subscription{}, fmt.Errorf("failed to get subscription by ID: %w", err)
	}

	// Check if the subscription is not deleted
	if deletedAt.Valid {
		// If deletedAt is valid, the subscription has been deleted
		return types.Subscription{}, ErrNotFound
	}

	return subscription, nil
}

// GetSubscriptionByCode retrieves a subscription by its code
func (s *Store) GetSubscriptionByCode(ctx context.Context, code string) (types.Subscription, error) {
	var subscription types.Subscription
	var deletedAt sql.NullTime

	err := s.getSubscriptionByCodeStmt.QueryRow(code).Scan(
		&subscription.SubscriptionID,
		&subscription.Code,
		&subscription.Plan,
		&subscription.Description,
		&subscription.Price,
		&subscription.DataLimit,
		&subscription.Validity,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&deletedAt,
	)

	// Check if the subscription is not deleted
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Subscription{}, ErrNotFound
		}
		return types.Subscription{}, fmt.Errorf("failed to get subscription by code: %w", err)
	}

	// Check if the subscription is not deleted
	if deletedAt.Valid {
		// If deletedAt is valid, the subscription has been deleted
		return types.Subscription{}, ErrNotFound
	}

	return subscription, nil
}

// Update the subscription
func (s *Store) UpdateSubscription(ctx context.Context, subscription types.Subscription) (types.Subscription, error) {
	err := s.updateSubscriptionStmt.QueryRow(
		subscription.Code,
		subscription.Plan,
		subscription.Description,
		subscription.Price,
		subscription.DataLimit,
		subscription.Validity,
		subscription.SubscriptionID,
	).Scan(
		&subscription.SubscriptionID,
		&subscription.Code,
		&subscription.Plan,
		&subscription.Description,
		&subscription.Price,
		&subscription.DataLimit,
		&subscription.Validity,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Subscription{}, ErrNotFound
		}
		return types.Subscription{}, fmt.Errorf("failed to update subscription: %w", err)
	}

	return subscription, nil
}

// Delete the subscription
func (s *Store) DeleteSubscription(ctx context.Context, id string) error {
	_, err := s.deleteSubscriptionStmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
