// middleware/session.go
package middleware

import (
	"auth_service_hi/config"
	"github.com/go-redis/redis/v8"
	"net/http"
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
		rdb.Expire(config.Ctx, token, 7*24*3600)
		next(w, r)
	}
}
