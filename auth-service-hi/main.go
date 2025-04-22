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
