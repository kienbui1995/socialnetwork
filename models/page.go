package models

// Page struct
type Page struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
	Description string `json:"description,omitempty"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	Email       string `json:"email"`
	Followers   int    `json:"followers"`
	Posts       int    `json:"posts"`
	CreatedAt   int    `json:"created_at,omitempty"`
	UpdatedAt   int    `json:"updated_at,omitempty"`
	Status      int    `json:"status"`
}

// PageObject struct
type PageObject struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
}
