package model

type Order struct {
	ID              string `json:"id"`
	PublicID        string `json:"public_id" bson:"public_id"`
	ProfilePublicID string `json:"profile_public_id" bson:"profile_public_id"`
	CarrierPublicID string `json:"carrier_public_id" bson:"carrier_public_id"`
	StorePublicID   string `json:"store_public_id" bson:"store_public_id"`
	Description     string `json:"description"`
	Quantity        string `json:"quantity"`
	Date            string `json:"date"`
	DeliveryDate    string `json:"delivery_date" bson:"delivery_date"`
	DepartureDate   string `json:"departure_date" bson:"departure_date"`
	State           int    `json:"state"`
	Price           float64
	Reference       string `json:"reference"`
	ClientPhone     string `json:"client_phone" bson:"client_phone"`
	ClientName      string `json:"client_name" bson:"client_name"`
	ProfileID       string `json:"profile_id" bson:"profile_id"`
	Detail          OrderDetail
	Experience      Experience
	ArrivalLocation Location `json:"arrival_location" bson:"arrival_location"`
	ExitLocation    Location `json:"exit_location" bson:"exit_location"`
}
type OrderDetail struct {
	amount      float64 `json:"amount" bson:"amount"`
	price       float64 `json:"price" bson:"price"`
	description string  `json:"description" bson:"description"`
}
type Experience struct {
	Score       int    `json:"score" bson:"score"`
	Date        string `json:"date" bson:"date"`
	Description string `json:"description" bson:"description"`
}
type Location struct {
	Latitude  *string `json:"latitude"`
	Longitude *string `json:"longitude"`
	Address   *string `json:"address"`
	Locality  *string `json:"locality"`
	Name      *string `json:"name"`
}
