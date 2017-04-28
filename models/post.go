package models

// CorePost struct to include info of core Post
type CorePost struct {
	PostID    int    `json:"id"` //Id
	Message   string `json:"message"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	Status    int    `json:"status,omitempty"`
	Privacy   int    `json:"privacy,omitempty"` // 1: public; 2: followers; 3: private
}

// Post struct
type Post struct {
	CorePost
	From   Profile `json:"from"`
	To     Profile `json:"to"`
	Image  string  `json:"image"`
	Action int     `json:"action"` // 1: Post; 2: Share
}
