package models

// Photo struct
type Photo struct {
	UserObject
	PhotoID   int    `json:"id"` //Id
	Message   string `json:"message"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	Status    int    `json:"status,omitempty"`
	Privacy   int    `json:"privacy,omitempty"` // 1: public; 2: followers; 3: private
	Likes     int    `json:"likes,omitempty"`
	Comments  int    `json:"comments,omitempty"`
	Shares    int    `json:"shares,omitempty"`
	IsLiked   bool   `json:"is_liked"`
}
