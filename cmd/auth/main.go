package main

import (
	"go-auth/config"
	authdb "go-auth/db"
	"go-auth/handlers"
	"go-auth/middleware"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file (optional)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()
	log.Println("Config loaded:", cfg)

	// Initialize database
	database, err := authdb.InitDB(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer authdb.Close(database)

	log.Println("Database connected successfully")

	// Create tables
	if err := authdb.CreateTables(database); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	log.Println("Tables created/verified")

	// Setup routes
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/health", handlers.HealthCheckHandler())
	mux.HandleFunc("/register", handlers.RegisterHandler(database, cfg.JWTSecret))
	mux.HandleFunc("/login", handlers.LoginHandler(database, cfg.JWTSecret))
	mux.HandleFunc("/refresh", handlers.RefreshTokenHandler(database, cfg.JWTSecret))

	// Protected routes (authentication required)
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	mux.Handle("/profile", authMiddleware(http.HandlerFunc(handlers.GetProfileHandler(database))))
	mux.Handle("/change-password", authMiddleware(http.HandlerFunc(handlers.ChangePasswordHandler(database))))

	// Start server
	log.Printf("Starting auth service on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}