package models

// UpdateUserRequest represents a partial update request for a user.
type UpdateUserRequest struct {
	Name         *string `json:"name"`     // Optional: Name field
	Email        *string `json:"email"`    // Optional: Email field
	PasswordHash *string `json:"password"` // Optional: Password field
}
