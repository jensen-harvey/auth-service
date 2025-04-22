// handlers/signup.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"

	"github.com/google/uuid"
	"herb_immortal/auth_service_hi/utils"
)

type SignUpRequest struct {
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
}

func SignUpHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		id := req.Role + "_" + uuid.New().String()

		_, err = db.Exec("INSERT INTO "+req.Role+"s (id, email, password_hash, created_at) VALUES ($1, $2, $3, $4)",
			id, req.Email, hash, time.Now())
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Stub for OTP logic here
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"user_id": id})
	}
}
