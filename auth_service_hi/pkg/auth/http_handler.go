package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/herb-immortal/auth_service_hi/pkg/models"
)

// HTTPHandler handles HTTP requests for authentication
type HTTPHandler struct {
	authService *AuthService
}

// NewHTTPHandler creates a new HTTP handler for authentication
func NewHTTPHandler(authService *AuthService) *HTTPHandler {
	return &HTTPHandler{authService: authService}
}

// EnableCORS adds CORS headers to responses
func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Call the next handler
		next(w, r)
	}
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// SignupHandler handles user registration
func (h *HTTPHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.SignupRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := h.authService.Signup(req)
	if err != nil {
		switch err {
		case ErrUserAlreadyExists:
			RespondWithError(w, http.StatusConflict, err.Error())
		case ErrInvalidRole:
			RespondWithError(w, http.StatusBadRequest, err.Error())
		default:
			RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		}
		return
	}

	RespondWithJSON(w, http.StatusCreated, user)
}

// LoginHandler authenticates users
func (h *HTTPHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	authResponse, err := h.authService.Login(req)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		case ErrInvalidOTP:
			RespondWithError(w, http.StatusUnauthorized, "Invalid OTP code")
		default:
			RespondWithError(w, http.StatusInternalServerError, "Error during login")
		}
		return
	}

	// Set session cookie for browser clients
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    authResponse.Token,
		Expires:  authResponse.ExpiresAt,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   true, // Set to true in production with HTTPS
	})

	RespondWithJSON(w, http.StatusOK, authResponse)
}

// Type to store user in context
type userContextKey string

const userKey userContextKey = "user"

// AuthMiddleware validates JWT tokens for protected routes
func (h *HTTPHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			// Try from cookie as fallback
			cookie, err := r.Cookie("session_token")
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Authorization token required")
				return
			}
			tokenString = cookie.Value
		} else {
			// Remove "Bearer " prefix if present
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]
			}
		}

		// Validate token
		user, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Store user in request context
		ctx := context.WithValue(r.Context(), userKey, user)
		next(w, r.WithContext(ctx))
	}
}

// WithUser adds the user to the request context
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) *models.User {
	if user, ok := ctx.Value(userKey).(*models.User); ok {
		return user
	}
	return nil
}

// SetupRoutes registers the authentication routes
func (h *HTTPHandler) SetupRoutes(mux *http.ServeMux) {
	// Apply CORS middleware to all routes
	mux.HandleFunc("/api/auth/signup", EnableCORS(h.SignupHandler))
	mux.HandleFunc("/api/auth/login", EnableCORS(h.LoginHandler))
	mux.HandleFunc("/api/auth/profile", EnableCORS(h.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		RespondWithJSON(w, http.StatusOK, user)
	})))
}