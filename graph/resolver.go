package graph

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/earqq/encargo-backend/db"
	"github.com/globalsign/mgo"
)

type Resolver struct {
	carriers *mgo.Collection
	orders   *mgo.Collection
	stores   *mgo.Collection
}

func New() Config {
	return Config{
		Resolvers: &Resolver{
			carriers: db.GetCollection("carriers"),
			orders:   db.GetCollection("orders"),
			stores:   db.GetCollection("stores"),
		},
	}
}
