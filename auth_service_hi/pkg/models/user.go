package models

import (
	"time"
)

// UserRole defines the possible roles a user can have
type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
	RoleHealer   UserRole = "healer"
	RoleVendor   UserRole = "vendor"
)

// User represents the basic user model that all user types will embed
type User struct {
	ID           string    `json:"id" db:"id"`                       // UUID with role prefix like "cust_123"
	Email        string    `json:"email" db:"email"`                 // Email address (unique)
	PasswordHash string    `json:"-" db:"password_hash"`             // Bcrypt hashed password
	MFASecret    string    `json:"-" db:"mfa_secret"`                // Encrypted MFA secret
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`   // Phone number
	Name         string    `json:"name" db:"name"`                   // User's name
	Role         UserRole  `json:"role" db:"role"`                   // User role (customer, admin, etc.)
	EmailVerified bool      `json:"email_verified" db:"email_verified"` // Whether email has been verified
	PhoneVerified bool      `json:"phone_verified" db:"phone_verified"` // Whether phone has been verified
	CreatedAt    time.Time `json:"created_at" db:"created_at"`       // Account creation timestamp
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`       // Account last update timestamp
}

// SignupRequest represents the data needed for signup
type SignupRequest struct {
	Name        string   `json:"name" binding:"required"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=8"`
	PhoneNumber string   `json:"phone_number" binding:"required"`
	Role        UserRole `json:"role" binding:"required"`
}

// LoginRequest represents the data needed for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	OTPCode  string `json:"otp_code,omitempty"` // Optional during initial login, required if 2FA is enabled
}

// AuthResponse represents the data returned after successful authentication
type AuthResponse struct {
	Token        string    `json:"token"`         // JWT token
	RefreshToken string    `json:"refresh_token"` // Refresh token (optional)
	ExpiresAt    time.Time `json:"expires_at"`    // Token expiration time
	User         User      `json:"user"`          // User information
}