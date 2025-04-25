// middleware/session.go
package middleware

import (
	"github.com/go-redis/redis/v8"
	"herb_immortal/auth_service_hi/config"
	"net/http"
	"time"
)

func SessionMiddleware(rdb *redis.Client, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		data, err := rdb.HGetAll(config.Ctx, token).Result()
		if err != nil || len(data) == 0 {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Optional: Extend session
		rdb.Expire(config.Ctx, token, 7*24*time.Hour)
		next(w, r)
	}
}
