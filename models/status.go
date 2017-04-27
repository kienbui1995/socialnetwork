package models

// UserStatus struct
type UserStatus struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userid,omitempty"`
	Message   string `json:"message"`
	CreatedAt int    `json:"created_at,omitempty"`
	UpdatedAt int    `json:"updated_at,omitempty"`
}
