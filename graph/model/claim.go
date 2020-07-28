package model

import jwt "github.com/dgrijalva/jwt-go"

type Claim struct {
	Username string `json:"username"`
	UserType string `json:"user_type" bson:"user_type"`
	jwt.StandardClaims
}
