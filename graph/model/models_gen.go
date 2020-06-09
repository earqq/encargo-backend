// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AddLocation struct {
	Latitude  *string `json:"latitude"`
	Longitude *string `json:"longitude"`
	Address   *string `json:"address"`
	Name      *string `json:"name"`
}

type CarrierStats struct {
	Orders  int     `json:"orders"`
	Ranking float64 `json:"ranking"`
}

type FilterOptions struct {
	Limit           int     `json:"limit"`
	PublicID        *string `json:"public_id"`
	State           *int    `json:"state"`
	State1          *int    `json:"state1"`
	State2          *int    `json:"state2"`
	CarrierPublicID *string `json:"carrier_public_id"`
	Search          *string `json:"search"`
}

type NewCarrier struct {
	StoreID      string  `json:"store_id"`
	Name         string  `json:"name"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	MessageToken *string `json:"message_token"`
	Phone        string  `json:"phone"`
}

type NewOrder struct {
	Description     string       `json:"description"`
	PublicID        string       `json:"public_id"`
	Reference       string       `json:"reference"`
	Price           float64      `json:"price"`
	ClientPhone     string       `json:"client_phone"`
	ClientName      string       `json:"client_name"`
	ExitLocation    *AddLocation `json:"exit_location"`
	ArrivalLocation *AddLocation `json:"arrival_location"`
}

type NewStore struct {
	PublicID string       `json:"public_id"`
	Name     string       `json:"name"`
	Phone    string       `json:"phone"`
	Location *AddLocation `json:"location"`
}

type UpdateCarrier struct {
	Name          *string `json:"name"`
	StateDelivery *int    `json:"state_delivery"`
	State         *bool   `json:"state"`
	Username      *string `json:"username"`
	Password      *string `json:"password"`
	MessageToken  *string `json:"message_token"`
	Phone         *string `json:"phone"`
}

type UpdateOrderInput struct {
	CarrierPublicID  *string `json:"carrier_public_id"`
	State            *int    `json:"state"`
	Score            *int    `json:"score"`
	ScoreDescription *string `json:"score_description"`
}
