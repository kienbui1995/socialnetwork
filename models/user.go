package models

// User stuct
type User struct {
	UserID        int    `json:"id" ` //Id
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	MiddleName    string `json:"middle_name,omitempty"`
	LastName      string `json:"last_name"`
	FullName      string `json:"full_name"`
	About         string `json:"about,omitempty"`
	Gender        string `json:"gender"`   //male or female
	Birthday      string `json:"birthday"` //This is a fixed format string, like DD/MM/YYYY
	Avatar        string `json:"avatar,omitempty"`
	Cover         string `json:"cover,omitempty"` //Direct URL for the user's cover photo image
	Status        int    `json:"status"`
	IsVertified   bool   `json:"is_vertified"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	CreatedAt     string `json:"created_at"`
	FacebookID    string `json:"facebook_id,omitempty"`
	FacebookToken string `json:"facebook_token,omitempty"`
	Posts         int    `json:"posts"`
	Followers     int    `json:"followers"`
	Followings    int    `json:"followings"`
}
