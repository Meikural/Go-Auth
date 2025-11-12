package main

import (
	"go-auth/config"
	seeder "go-auth/db/Seeder"
	authdb "go-auth/db/database"
	"go-auth/handlers"
	"go-auth/handlers/admin"
	"go-auth/handlers/auth"
	"go-auth/handlers/user"
	"go-auth/middleware"
	authmiddle "go-auth/middleware/auth"
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

	// Seed roles
	if err := seeder.SeedRoles(database, cfg.Roles); err != nil {
		log.Fatalf("Failed to seed roles: %v", err)
	}

	log.Println("Roles seeded successfully")

	// Seed super admin
	superAdminRole := cfg.Roles[0] // First role is Super Admin
	if err := seeder.SeedSuperAdmin(database, cfg.SuperAdminEmail, cfg.SuperAdminPassword, superAdminRole); err != nil {
		log.Fatalf("Failed to seed super admin: %v", err)
	}

	log.Println("Super admin seeded successfully")

	// Setup routes
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/health", handlers.HealthCheckHandler())
	mux.HandleFunc("/register", auth.RegisterHandler(database, cfg.JWTSecret, cfg.DefaultRegistrationRole))
	mux.HandleFunc("/login", auth.LoginHandler(database, cfg.JWTSecret))
	mux.HandleFunc("/refresh", auth.RefreshTokenHandler(database, cfg.JWTSecret))

	// Protected routes (authentication required)
	authMiddleware := authmiddle.AuthMiddleware(cfg.JWTSecret)
	mux.Handle("/profile", authMiddleware(http.HandlerFunc(user.GetProfileHandler(database))))
	mux.Handle("/change-password", authMiddleware(http.HandlerFunc(auth.ChangePasswordHandler(database))))

	// Admin routes (authentication + role required)
	roleMiddleware := middleware.RequireRole(superAdminRole)
	
	// Get all users
	mux.Handle("/admin/users", authMiddleware(roleMiddleware(http.HandlerFunc(admin.GetAllUsersHandler(database)))))
	
	// Get specific user: GET /admin/users/get/{uuid}
	mux.Handle("/admin/users/get/", authMiddleware(roleMiddleware(http.HandlerFunc(admin.GetUserHandler(database)))))
	
	// Create user: POST /admin/users/create
	mux.Handle("/admin/users/create", authMiddleware(roleMiddleware(http.HandlerFunc(admin.CreateUserHandler(database, cfg.Roles)))))
	
	// Update user: PATCH /admin/users/update/{uuid}
	mux.Handle("/admin/users/update/", authMiddleware(roleMiddleware(http.HandlerFunc(admin.UpdateUserHandler(database)))))
	
	// Delete user: DELETE /admin/users/delete/{uuid}
	mux.Handle("/admin/users/delete/", authMiddleware(roleMiddleware(http.HandlerFunc(admin.DeleteUserHandler(database)))))
	
	// Update user role: PUT /admin/users/role/{uuid}
	mux.Handle("/admin/users/role/", authMiddleware(roleMiddleware(http.HandlerFunc(admin.UpdateUserRoleHandler(database, cfg.Roles)))))

	// Start server
	log.Printf("Starting auth service on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}