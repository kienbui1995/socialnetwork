package models

// User stuct
type User struct {
	Data map[string]interface{}
	// UserID      int    `form:"userid" json:"userid" binding:"required"` //Id
	// Username    string `form:"username" json:"username" binding:"username"`
	// Password    string `form:"password" json:"password" binding:"password"`
	// Email       string `form:"email" json:"email" binding:"email"`
	// FirstName   string `form:"firstname" json:"firstname" binding:"firstname"`
	// MiddleName  string `form:"middlename" json:"middlename" binding:"middlename"`
	// LastName    string `form:"lastname" json:"lastname" binding:"lastname"`
	// FullName    string `form:"fullname" json:"fullname" binding:"fullname"`
	// About       string `form:"about" json:"about" binding:"about"`
	// Gender      string `form:"gender" json:"gender" binding:"gender"`
	// Birthday    string `form:"birthday" json:"birthday" binding:"birthday"` //This is a fixed format string, like MM/DD/YYYY
	// Avatar      string `form:"avatar" json:"avatar" binding:"avatar	"`
	// CoverPhoto  string `form:"coverphoto" json:"coverphoto" binding:"coverphoto"` //Direct URL for the user's cover photo image
	// Status      int    `form:"status" json:"status" binding:"status"`
	// IsVertified bool   `form:"isvertified" json:"isvertified" binding:"isvertified"`
	// UpdatedTime string `form:"updatetime" json:"updatetime" binding:"updatetime"`
	// CreatedTime string `form:"createtime" json:"createtime" binding:"createtime"`
}
