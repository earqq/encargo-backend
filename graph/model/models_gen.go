// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AddLocation struct {
	Latitude  *string `json:"latitude"`
	Longitude *string `json:"longitude"`
	Address   *string `json:"address"`
	Reference *string `json:"reference"`
}

type CarrierStats struct {
	Orders  int     `json:"orders"`
	Ranking float64 `json:"ranking"`
}

type FilterOptions struct {
	Limit     int       `json:"limit"`
	State     *int      `json:"state"`
	State1    *int      `json:"state1"`
	State2    *int      `json:"state2"`
	CarrierID *string   `json:"carrier_id"`
	Search    *string   `json:"search"`
	Ids       []*string `json:"ids"`
}

type NewCarrier struct {
	StoreID      *string `json:"store_id"`
	Name         string  `json:"name"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	MessageToken *string `json:"message_token"`
	Phone        string  `json:"phone"`
}

type NewOrder struct {
	Price           float64           `json:"price"`
	DeliveryPrice   *float64          `json:"delivery_price"`
	ClientPhone     string            `json:"client_phone"`
	ClientName      string            `json:"client_name"`
	ArrivalLocation *AddLocation      `json:"arrival_location"`
	Detail          []*NewOrderDetail `json:"detail"`
}

type NewOrderDetail struct {
	Amount      *float64 `json:"amount"`
	Price       *float64 `json:"price"`
	Description *string  `json:"description"`
}

type NewStore struct {
	Name       string       `json:"name"`
	Phone      string       `json:"phone"`
	Ruc        *string      `json:"ruc"`
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	FirebaseID *string      `json:"firebaseID"`
	Location   *AddLocation `json:"location"`
}

type UpdateCarrier struct {
	Name          *string `json:"name"`
	StateDelivery *int    `json:"state_delivery"`
	Global        *bool   `json:"global"`
	Password      *string `json:"password"`
	MessageToken  *string `json:"message_token"`
}

type UpdateCarrierLocation struct {
	ActualLocation *AddLocation `json:"actual_location"`
}

type UpdateOrder struct {
	CarrierID        *string `json:"carrier_id"`
	State            *int    `json:"state"`
	Score            *int    `json:"score"`
	ScoreDescription *string `json:"score_description"`
}

type UpdateStore struct {
	Name     *string      `json:"name"`
	Location *AddLocation `json:"location"`
}
