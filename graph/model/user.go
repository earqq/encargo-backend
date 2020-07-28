package model

type User struct {
	ID        string `json:"id" bson:"_id"`
	UserType        string `json:"user_type" bson:"user_type"`
	Username        string `json:"username" bson:"username"`
}
