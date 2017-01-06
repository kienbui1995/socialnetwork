package models

// User stuct
type User struct {
	UserID   int    `form:"userid" json:"userid" binding:"required"` //Id
	Username string `form:"username" json:"username" binding:"username"`
	Password string `form:"password" json:"password" binding:"password"`
	Email    string `form:"email" json:"email" binding:"email"`
	Status   int    `form:"status" json:"status" binding:"status"`
}
