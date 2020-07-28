package model

type Store struct {
	ID         string   `json:"_id" bson:"_id"`
	FirebaseID string   `json:"firebase_id" bson:"firebase_id"`
	Name       string   `json:"name"`
	Username   string   `json:"username"`
	Ruc   string   `json:"ruc"`
	Password   string   `json:"password"`
	Phone      string   `json:"phone"`
	Token         string     `json:"token"`
	Location   Location `json:"location" bson:"location"`
}
