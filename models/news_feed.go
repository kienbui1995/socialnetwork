package models

// News struct
type News struct {
	UserObject
	ID        int    `json:"id"`
	Message   string `json:"message"`
	Summary   bool   `json:"summary,omitempty"`
	Photo     string `json:"photo,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	Privacy   int    `json:"privacy,omitempty"` // 1: public; 2: followers; 3: private
	Status    int    `json:"status,omitempty"`
	Likes     int    `json:"likes,omitempty"`
	Comments  int    `json:"comments,omitempty"`
	Shares    int    `json:"shares,omitempty"`
	IsLiked   bool   `json:"is_liked"`
	CanEdit   bool   `json:"can_edit"`
	CanDelete bool   `json:"can_delete"`
}
