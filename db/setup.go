package db

import (
	"fmt"
	"github.com/globalsign/mgo"
	"os"
)
var session *mgo.Session
var db *mgo.Database
const defaultDB = "encargo"
func ConnectDB() {
	session, err := mgo.Dial("mongodb://localhost:27017,localhost:27018")
	if err != nil {
		fmt.Println(err)
	}
	session.SetMode(mgo.Monotonic, true)
	databaseName := os.Getenv("DATABASE")
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
