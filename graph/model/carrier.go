package model

type Carrier struct {
	ID             string `json:"_id" bson:"_id"`
	StoreID        string `json:"store_id" bson:"store_id"`
	Name           string
	Birthdate      string
	StateDelivery  int `json:"state_delivery" bson:"state_delivery"`
	Username       string
	Password       string
	CurrentOrderID string `json:"current_order_id" bson:"current_order_id"`
	MessageToken   string `json:"message_token" bson:"message_token"`
	Phone          string
}
