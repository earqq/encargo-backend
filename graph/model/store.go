package model

type Store struct {
	ID         string   `json:"_id" bson:"_id"`
	FirebaseID string   `json:"firebase_id" bson:"firebase_id"`
	Name       string   `json:"name"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	Phone      string   `json:"phone"`
	Location   Location `json:"location" bson:"location"`
}
