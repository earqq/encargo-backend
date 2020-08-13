package model

type Carrier struct {
	ID             string `json:"_id" bson:"_id"`
	StoreID        string `json:"store_id" bson:"store_id"`
	Global         bool   `json:"global" bson:"global"`
	Name           string
	Birthdate      string
	StateDelivery  int  `json:"state_delivery" bson:"state_delivery"`
	State          bool `json:"state" bson:"state"`
	Username       string
	Password       string
	Token          string `json:"token"`
	CurrentOrderID string `json:"current_order_id" bson:"current_order_id"`
	MessageToken   string `json:"message_token" bson:"message_token"`
	Phone          string
}
