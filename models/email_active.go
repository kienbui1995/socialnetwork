package models

//EmailActive struct
type EmailActive struct {
	UserID     int    `json:"id"`
	Email      string `json:"email"`
	ActiveCode string `json:"active_code"`
}
