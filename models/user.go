package models

// User stuct
type User struct {
	UserID        int    `json:"id"` //Id
	Username      string `json:"username"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email" binding:"required"`
	FirstName     string `json:"first_name,omitempty"`
	MiddleName    string `json:"middle_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	FullName      string `json:"full_name,omitempty"`
	About         string `json:"about,omitempty"`
	Gender        string `json:"gender,omitempty"`       //male or female
	Birthday      string `json:"birthday,omitempty"`     //This is a fixed format string, like DD/MM/YYYY
	Avatar        string `json:"avatar,omitempty"`       //Direct URL for the user's avatar image
	Cover         string `json:"cover,omitempty"`        //Direct URL for the user's cover image
	Status        int    `json:"status,omitempty"`       //1: Active; 0: DeActive
	IsVertified   bool   `json:"is_vertified,omitempty"` // true or false
	UpdatedAt     string `json:"updated_at,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	FacebookID    string `json:"facebook_id,omitempty"`
	FacebookToken string `json:"facebook_token" valid:"omitempty"`
	Posts         int    `json:"posts"`
	Followers     int    `json:"followers"`
	Followings    int    `json:"followings"`
	EmailActive   string `json:"email_active,omitempty"`
}
