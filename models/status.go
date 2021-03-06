package models

// UserStatus struct
type UserStatus struct {
	UserObject
	ID        int    `json:"id"`
	Message   string `json:"message"`
	CreatedAt int    `json:"created_at,omitempty"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	Privacy   int    `json:"privacy,omitempty"` // 1: public; 2: followers; 3: private
	Status    int    `json:"status,omitempty"`
	Likes     int    `json:"likes,omitempty"`
	Comments  int    `json:"comments,omitempty"`
	Shares    int    `json:"shares,omitempty"`
	IsLiked   bool   `json:"is_liked"`
}
