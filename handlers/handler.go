package handlers

import (
	"net/http"
)

// DataResponse represents the response from protected endpoint
type DataResponse struct {
	Message  string `json:"message"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Data     string `json:"data"`
}

// GetDataHandler is a protected endpoint that returns data
func GetDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			RespondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Get user info from headers (set by middleware)
		userID := r.Header.Get("X-User-ID")
		username := r.Header.Get("X-Username")
		email := r.Header.Get("X-Email")
		role := r.Header.Get("X-Role")

		response := DataResponse{
			Message:  "success",
			UserID:   userID,
			Username: username,
			Email:    email,
			Role:     role,
			Data:     "This is protected data accessible only with valid access token",
		}

		RespondJSON(w, http.StatusOK, response)
	}
}
