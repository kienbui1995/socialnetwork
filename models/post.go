package models

// Post struct to include info of Post
type Post struct {
	ID       int64  `form:"id" json:"id" binding:"required"` //Id
	Content  string `form:"username" json:"username" binding:"username"`
	Type     string `form:"password" json:"password" binding:"password"`
	Status   string `form:"email" json:"email" binding:"email"`
	IsHidden bool   `form:"is_hidden" json:"is_hidden" binding:"is_hidden"`
}
