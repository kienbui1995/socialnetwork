package models

// UserComment struct
type UserComment struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userid"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Avatar    string `json:"avatar"`
	Message   string `json:"message"`
	CreatedAt int    `json:"created_at,omitempty"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	Status    int    `json:"status,omitempty"`
}
