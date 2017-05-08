package models

// Post struct to include info of core Post
type Post struct {
	UserObject
	PostID    int    `json:"id"` //Id
	Message   string `json:"message"`
	Photo     string `json:"photo,omitempty"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	Status    int    `json:"status,omitempty"`
	Privacy   int    `json:"privacy,omitempty"` // 1: public; 2: followers; 3: private
	Likes     int    `json:"likes,omitempty"`
	Comments  int    `json:"comments,omitempty"`
	Shares    int    `json:"shares,omitempty"`
	IsLiked   bool   `json:"is_liked"`
	CanEdit   bool   `json:"can_edit"`
	CanDelete bool   `json:"can_delete"`
}
