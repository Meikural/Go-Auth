package admin

import (
	"database/sql"
	queries "go-auth/db/Queries"
	"go-auth/handlers"
	"go-auth/models"
	"net/http"
)

// GetAllUsersHandler returns all users (Super Admin only)
func GetAllUsersHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			
			
		}

		// Get all users
		users, err := queries.GetAllUsers(database)
		if err != nil {
			handlers.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get users"})
			return
		}

		response := models.GetAllUsersResponse{
			Total: len(users),
			Users: users,
		}

		handlers.RespondJSON(w, http.StatusOK, response)
	}
}