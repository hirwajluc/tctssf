package models

// SessionData represents user session information
type SessionData struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}

// LoginRequest for authentication
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse for successful authentication
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}