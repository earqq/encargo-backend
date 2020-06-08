package model

type Store struct {
	ID       string `json:"id"`
	PublicID string `json:"public_id" bson:"public_id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	// Location Location `json:"location" bson:"location"`
}
