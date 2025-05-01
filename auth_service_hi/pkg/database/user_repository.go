package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/herb-immortal/auth_service_hi/pkg/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
	INSERT INTO users (id, email, password_hash, mfa_secret, phone_number, name, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.MFASecret,
		user.PhoneNumber,
		user.Name,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
	SELECT id, email, password_hash, mfa_secret, phone_number, name, role, email_verified, phone_verified, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.MFASecret,
		&user.PhoneNumber,
		&user.Name,
		&user.Role,
		&user.EmailVerified,
		&user.PhoneVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	query := `
	SELECT id, email, password_hash, mfa_secret, phone_number, name, role, email_verified, phone_verified, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.MFASecret,
		&user.PhoneNumber,
		&user.Name,
		&user.Role,
		&user.EmailVerified,
		&user.PhoneVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// UpdateVerificationStatus updates the verification status of a user's email or phone
func (r *UserRepository) UpdateVerificationStatus(userID string, emailVerified, phoneVerified bool) error {
	query := `
	UPDATE users
	SET email_verified = $1, phone_verified = $2, updated_at = $3
	WHERE id = $4
	`

	_, err := r.db.Exec(query, emailVerified, phoneVerified, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update verification status: %w", err)
	}

	return nil
}

// SaveSession stores a session in the database
func (r *UserRepository) SaveSession(sessionID, userID string, expiresAt time.Time) error {
	query := `
	INSERT INTO sessions (id, user_id, expires_at)
	VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, sessionID, userID, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// GetSession retrieves a session by ID
func (r *UserRepository) GetSession(sessionID string) (string, time.Time, error) {
	query := `
	SELECT user_id, expires_at
	FROM sessions
	WHERE id = $1
	`

	var userID string
	var expiresAt time.Time

	err := r.db.QueryRow(query, sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", time.Time{}, nil
		}
		return "", time.Time{}, fmt.Errorf("failed to get session: %w", err)
	}

	return userID, expiresAt, nil
}

// DeleteSession removes a session
func (r *UserRepository) DeleteSession(sessionID string) error {
	query := `
	DELETE FROM sessions
	WHERE id = $1
	`

	_, err := r.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}