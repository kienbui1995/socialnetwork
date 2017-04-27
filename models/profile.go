package models

//Profile interface include User, Page, Group
type Profile struct {
	ID     int `json:"id"`
	Object interface{}
}
