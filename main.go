// main.go
package main

import (
	"log"
	"net/http"

	"herb_immortal/auth_service_hi/config"
	"herb_immortal/auth_service_hi/handlers"
)

func main() {
	config.LoadEnv()
	db := config.InitDB()
	rdb := config.InitRedis()

	http.HandleFunc("/signup", handlers.SignUpHandler(db, rdb))
	http.HandleFunc("/login", handlers.LoginHandler(db, rdb))

	log.Println("[INFO] Auth service started on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// config/config.go
package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

// config/db.go
package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func InitDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	return db
}

// config/redis.go
package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var Ctx = context.Background()

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})
	return rdb
}

// handlers/signup.go
package handlers

import (
	"database/sql"
	"encoding/json"
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

		_, err = db.Exec("INSERT INTO " + req.Role + "s (id, email, password_hash, created_at) VALUES ($1, $2, $3, $4)",
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

// utils/hashing.go
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
