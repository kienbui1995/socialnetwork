package models

// Post struct to include info of Post
type Post struct {
	PostID      int    `form:"id" json:"id" binding:"required"` //Id
	Content     string `form:"content" json:"content" binding:"content"`
	Image       string `form:"image" json:"image" binding:"image"`
	UpdatedTime string `form:"updatedtime" json:"updatedtime" binding:"updatedtime"`
	CreatedTime string `form:"createtime" json:"createtime" binding:"createtime"`
	Status      int    `form:"status" json:"status" binding:"status"`
}
