package models

// User stuct
type User struct {
	Account
	FirstName     string `json:"first_name,omitempty"`
	MiddleName    string `json:"middle_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	FullName      string `json:"full_name,omitempty"`
	About         string `json:"about,omitempty"`
	Gender        string `json:"gender,omitempty"`   //male or female
	Birthday      string `json:"birthday,omitempty"` //This is a fixed format string, like DD/MM/YYYY
	Avatar        string `json:"avatar,omitempty"`   //Direct URL for the user's avatar image
	Cover         string `json:"cover,omitempty"`    //Direct URL for the user's cover image
	Status        int    `json:"status,omitempty"`   //1: Active; 0: DeActive
	UpdatedAt     string `json:"updated_at,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	FacebookToken string `json:"facebook_token,omitempty"`
	Posts         int    `json:"posts"`
	Followers     int    `json:"followers"`
	Followings    int    `json:"followings"`
}
