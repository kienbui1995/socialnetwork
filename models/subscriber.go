package models

//Subscriber struct
type Subscriber struct {
	SubscriberID int `json:"id"`
	FromID       int `json:"from_id"`
	ToID         int `json:"to_id"`
}
