package model

type Carrier struct {
	ID             string `json:"_id" bson:"_id"`
	PublicID       string `json:"public_id" bson:"public_id"`
	StorePublicID  string `json:"store_public_id" bson:"store_public_id"`
	Name           string
	Birthdate      string
	StateDelivery  int `json:"state_delivery" bson:"state_delivery"`
	Username       string
	Password       string
	CurrentOrderID string `json:"current_order_id" bson:"current_order_id"`
	MessageToken   string `json:"message_token" bson:"message_token"`
	Phone          string
}
