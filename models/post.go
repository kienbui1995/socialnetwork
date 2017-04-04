package models

// Post struct to include info of Post
type Post struct {
	PostID      int    `json:"id"` //Id
	Content     string `json:"content"`
	Image       string `json:"image"`
	UpdatedTime string `json:"updated_time"`
	CreatedTime string `json:"created_time"`
	Status      int    `json:"status"`
}
