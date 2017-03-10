package models

// Post struct to include info of Post
type Post struct {
	PostID       int    `form:"id" json:"id" binding:"required"` //Id
	Message      string `form:"message" json:"message" binding:"message"`
	Type         string `form:"password" json:"password" binding:"password"`
	Status       int    `form:"status" json:"status" binding:"status"`
	IsHidden     bool   `form:"ishidden" json:"ishidden" binding:"ishidden"`
	UpdatedTime  string `form:"updatedtime" json:"updatedtime" binding:"updatedtime"`
	CreatedTime  string `form:"createtime" json:"createtime" binding:"createtime"`
	CreateUserID int    `form:"createuserid" json:"createuserid" binding:"createuserid"`
}
