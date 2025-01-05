package user

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
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already exists")
	ErrUpdateFailed   = errors.New("failed to update user")
	ErrDeleteFailed   = errors.New("failed to delete user")
	ErrUserDeleted    = errors.New("user is deleted, please contact support")
)

type Store struct {
	db *sql.DB
	// Pre-prepared statements
	createUserStmt     *sql.Stmt
	getUserByEmailStmt *sql.Stmt
	getUserByIDStmt    *sql.Stmt
	updateUserStmt     *sql.Stmt
	deleteUserStmt     *sql.Stmt
	changePasswordStmt *sql.Stmt
}

// Queries
const (
	createUserQuery = `
		INSERT INTO users (
			firstName, 
			lastName, 
			email, 
			password
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	getUserByEmailQuery = `
		SELECT 
			id, 
			firstName, 
			lastName, 
			email, 
			password, 
			created_at, 
			updated_at,
			deleted_at
		FROM users 
		WHERE email = $1`

	getUserByIDQuery = `
		SELECT 
			id, 
			firstName, 
			lastName, 
			email, 
			password, 
			created_at, 
			updated_at,
			deleted_at
		FROM users 
		WHERE id = $1`

	updateUserQuery = `
		UPDATE users
		SET
			firstName = $1,
			lastName = $2,
			email = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`

	deleteUserQuerySoft = `
	UPDATE users 
	SET 
		deleted_at = NOW() 
	WHERE id = $1`

	changePasswordStmt = `
		UPDATE users
		SET
			password = $1,
			updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at`
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

	s.createUserStmt, err = s.db.Prepare(createUserQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare create user statement: %w", err)
	}

	s.getUserByEmailStmt, err = s.db.Prepare(getUserByEmailQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare get user by email statement: %w", err)
	}

	s.getUserByIDStmt, err = s.db.Prepare(getUserByIDQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare get user by ID statement: %w", err)
	}

	s.updateUserStmt, err = s.db.Prepare(updateUserQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare update user statement: %w", err)
	}

	s.deleteUserStmt, err = s.db.Prepare(deleteUserQuerySoft)
	if err != nil {
		return fmt.Errorf("failed to prepare delete user statement: %w", err)
	}

	s.changePasswordStmt, err = s.db.Prepare(changePasswordStmt)
	if err != nil {
		return fmt.Errorf("failed to prepare change password statement: %w", err)
	}

	return nil
}

// Close releases all prepared statements
func (s *Store) Close() error {
	var errs []error

	if err := s.createUserStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close create user stmt: %w", err))
	}
	if err := s.getUserByEmailStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close get user by email stmt: %w", err))
	}
	if err := s.getUserByIDStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close get user by ID stmt: %w", err))
	}
	if err := s.updateUserStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close update user stmt: %w", err))
	}
	if err := s.deleteUserStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close delete user stmt: %w", err))
	}
	if err := s.changePasswordStmt.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close change password stmt: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close statements: %v", errs)
	}

	return nil
}

// CreateUser creates a new user and returns the created user with ID and timestamps
func (s *Store) CreateUser(ctx context.Context, user types.User) (types.User, error) {
	err := s.createUserStmt.QueryRowContext(
		ctx,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if isPgDuplicateKeyError(err) {
			return types.User{}, ErrDuplicateEmail
		}
		return types.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (s *Store) GetUserByEmail(ctx context.Context, email string) (types.User, error) {
	var user types.User
	var deletedAt sql.NullTime

	err := s.getUserByEmailStmt.QueryRowContext(ctx, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, ErrUserNotFound
		}
		return types.User{}, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}

	// Check if user is not deleted
	if deletedAt.Valid {
		// If deletedAt is valid, the user has been deleted
		return types.User{}, ErrUserDeleted
	}

	// Convert sql.NullTime to time.Time if it's not NULL
	user.DeletedAt = deletedAt.Time

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *Store) GetUserByID(ctx context.Context, id string) (types.User, error) {
	var user types.User
	var deletedAt sql.NullTime

	err := s.getUserByIDStmt.QueryRowContext(ctx, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, ErrUserNotFound
		}
		return types.User{}, fmt.Errorf("failed to get user by id %s: %w", id, err)
	}

	// Check if user is not deleted
	if deletedAt.Valid {
		// If deletedAt is valid, the user has been deleted
		return types.User{}, ErrUserDeleted
	}

	// Convert sql.NullTime to time.Time if it's not NULL
	user.DeletedAt = deletedAt.Time

	return user, nil
}

// UpdateUser updates a user's details
func (s *Store) UpdateUser(ctx context.Context, user types.User) (types.User, error) {
	err := s.updateUserStmt.QueryRowContext(
		ctx,
		user.FirstName,
		user.LastName,
		user.Email,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, ErrUserNotFound
		}
		return types.User{}, ErrUpdateFailed
	}

	return user, nil
}

// DeleteUser removes a user from the database
func (s *Store) DeleteUser(ctx context.Context, id string) error {
	_, err := s.deleteUserStmt.ExecContext(ctx, id)
	if err != nil {
		return ErrDeleteFailed
	}
	return nil
}

// ChangePassword updates a user's password
func (s *Store) ChangePassword(ctx context.Context, user types.User) error {
	_, err := s.changePasswordStmt.ExecContext(ctx, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	return nil
}

// isPgDuplicateKeyError checks if the error is a PostgreSQL unique violation
func isPgDuplicateKeyError(err error) bool {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return pqErr.Code == "23505" // unique_violation error code
}
