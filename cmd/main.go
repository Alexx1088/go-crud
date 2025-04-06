package main

import (
	"fmt"
	"go-crud/internal/config"
	"go-crud/internal/handlers"
	"go-crud/internal/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Access environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	//Load JWT secret key
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	utils.SetJWTSecret(jwtSecret)

	// Construct connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println("Database connection string:", connStr)

	// Initialize database connection
	db := config.InitDB(connStr)
	defer db.Close()

	// Initialize router
	router := mux.NewRouter()

	// Register routes
	handlers.RegisterUserRoutes(router, db, []byte(jwtSecret))

	handlers.RegisterAuthRoutes(router, db)

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
