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
