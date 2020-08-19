package db

import (
	"fmt"
	"log"
	"os"

	"github.com/globalsign/mgo"
	"github.com/joho/godotenv"
)

var session *mgo.Session
var db *mgo.Database

const defaultDB = "encargo"

func ConnectDB() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	session, err := mgo.Dial("mongodb://localhost:27017,localhost:27018")
	if err != nil {
		fmt.Println(err)
	}
	session.SetMode(mgo.Monotonic, true)
	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = defaultDB
	}
	db = session.DB(databaseName)
}
func GetCollection(collection string) *mgo.Collection {
	return db.C(collection)
}
func CloseSession() {
	session.Close()
}
