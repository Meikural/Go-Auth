package models

// UpdateRoleRequest is the payload for updating user role
type UpdateRoleRequest struct {
	Role string `json:"role"`
}

// GetAllUsersResponse is the response for getting all users
type GetAllUsersResponse struct {
	Total int     `json:"total"`
	Users []*User `json:"users"`
}

// UpdateRoleResponse is the response for updating user role
type UpdateRoleResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user"`
}