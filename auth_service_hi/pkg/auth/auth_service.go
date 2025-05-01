package auth

import (
	"errors"
	"time"

	"github.com/herb-immortal/auth_service_hi/pkg/database"
	"github.com/herb-immortal/auth_service_hi/pkg/models"
	"github.com/herb-immortal/auth_service_hi/pkg/utils"
)

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserAlreadyExists   = errors.New("user with this email already exists")
	ErrInvalidRole         = errors.New("invalid role")
	ErrInvalidOTP          = errors.New("invalid OTP code")
	ErrVerificationRequired = errors.New("email or phone verification required")
	ErrInvalidSession      = errors.New("invalid or expired session")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo     *database.UserRepository
	tokenManager *utils.TokenManager
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *database.UserRepository, tokenManager *utils.TokenManager) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

// Signup registers a new user
func (s *AuthService) Signup(req models.SignupRequest) (*models.User, error) {
	// Check if user with this email already exists
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Validate role
	if req.Role != models.RoleCustomer && 
	   req.Role != models.RoleAdmin && 
	   req.Role != models.RoleHealer && 
	   req.Role != models.RoleVendor {
		return nil, ErrInvalidRole
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user with role prefix in ID
	now := time.Now()
	user := &models.User{
		ID:           utils.GenerateUUID(req.Role),
		Email:        req.Email,
		PasswordHash: passwordHash,
		MFASecret:    "", // Would be generated and encrypted in a full implementation
		PhoneNumber:  req.PhoneNumber,
		Name:         req.Name,
		Role:         req.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save user to database
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// In a real implementation, you would send verification OTPs here
	// and initialize MFA if required

	return user, nil
}

// Login authenticates a user and returns a session token
func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// In a real implementation, you would verify OTP/2FA here
	// For simplicity, we'll skip that in this basic implementation
	// if user.MFASecret != "" && req.OTPCode == "" {
	//     return nil, errors.New("OTP code required")
	// }
	// if user.MFASecret != "" {
	//     // Verify OTP code
	//     valid := verifyOTP(user.MFASecret, req.OTPCode)
	//     if !valid {
	//         return nil, ErrInvalidOTP
	//     }
	// }

	// Generate JWT token
	token, expiresAt, err := s.tokenManager.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// Generate session ID (for database/Redis storage)
	sessionID := utils.GenerateUUID(models.UserRole("session"))

	// Store session in database
	err = s.userRepo.SaveSession(sessionID, user.ID, expiresAt)
	if err != nil {
		return nil, err
	}

	// Return response with token and user info
	return &models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.User{
			ID:           user.ID,
			Email:        user.Email,
			PhoneNumber:  user.PhoneNumber,
			Name:         user.Name,
			Role:         user.Role,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		},
	}, nil
}

// ValidateSession checks if a session is valid
func (s *AuthService) ValidateSession(sessionID string) (*models.User, error) {
	// Get session from database
	userID, expiresAt, err := s.userRepo.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	
	// Check if session exists and hasn't expired
	if userID == "" || expiresAt.Before(time.Now()) {
		return nil, ErrInvalidSession
	}
	
	// Get user by ID
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidSession
	}
	
	// Extend session validity (in a real implementation with Redis)
	// Here you would update the TTL in Redis
	
	return user, nil
}

// ValidateToken validates a JWT token and returns the associated user
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// Verify JWT token
	claims, err := s.tokenManager.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	// Get user by ID
	user, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidSession
	}
	
	return user, nil
}