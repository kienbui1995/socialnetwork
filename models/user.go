package models

// User stuct
type User struct {
	UserID      int    `json:"id" ` //Id
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	FullName    string `json:"full_name"`
	About       string `json:"about"`
	Gender      string `json:"gender"`
	Birthday    string `json:"birthday"` //This is a fixed format string, like MM/DD/YYYY
	Avatar      string `json:"avatar"`
	CoverPhoto  string `json:"cover_photo"` //Direct URL for the user's cover photo image
	Status      int    `json:"status"`
	IsVertified bool   `json:"is_vertified"`
	UpdatedTime string `json:"updated_time"`
	CreatedTime string `json:"created_time"`
}
