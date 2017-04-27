package models

// CorePost struct to include info of core Post
type CorePost struct {
	PostID    int    `json:"id"` //Id
	Message   string `json:"message"`
	UpdatedAt int    `json:"updated_at,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
	Status    int    `json:"status,omitempty"`
}

// Post struct
type Post struct {
	CorePost
	From Profile `json:"from"`
	To   Profile `json:"to"`
}
