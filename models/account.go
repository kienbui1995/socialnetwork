package models

//Account struct
type Account struct {
	UserID      int    `json:"id"` //Id
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	FacebookID  string `json:"facebook_id,omitempty"`
	IsVertified bool   `json:"is_vertified"` // true or false
}
