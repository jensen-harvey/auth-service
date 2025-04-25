// handlers/login.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"herb_immortal/auth_service_hi/config"
	"herb_immortal/auth_service_hi/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func LoginHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var id, hash string
		err := db.QueryRow("SELECT id, password_hash FROM "+req.Role+"s WHERE email=$1", req.Email).Scan(&id, &hash)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if err := utils.ComparePassword(hash, req.Password); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		sessionID := utils.GenerateUUID()
		expires := time.Now().Add(7 * 24 * time.Hour)
		rdb.HSet(config.Ctx, sessionID, map[string]interface{}{
			"user_id":    id,
			"role":       req.Role,
			"expires_at": expires.Unix(),
		})
		rdb.Expire(config.Ctx, sessionID, 7*24*time.Hour)

		token, _ := utils.GenerateJWT(id, req.Role)

		json.NewEncoder(w).Encode(map[string]string{
			"session_token": sessionID,
			"jwt":           token,
		})
	}
}
