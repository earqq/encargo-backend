package graph

//go:generate go run github.com/99designs/gqlgen
import (
	"math/rand"
	"sync"

	"github.com/earqq/encargo-backend/db"
	"github.com/earqq/encargo-backend/graph/generated"
	"github.com/earqq/encargo-backend/graph/model"
	"github.com/globalsign/mgo"
	"golang.org/x/crypto/bcrypt"
)

func New() generated.Config {
	return generated.Config{
		Resolvers: &Resolver{
			carriers:                     db.GetCollection("carriers"),
			orders:                       db.GetCollection("orders"),
			stores:                       db.GetCollection("stores"),
			storeOrdersTopics:            map[string]*StoreOrdersTopic{},
			orderTopics:                  map[string]*OrderTopic{},
			carrierTopics:                map[string]*CarrierTopic{},
			storeCarriersTopics:          map[string]*StoreCarriersTopic{},
			carrierLocationTopics:        map[string]*CarrierLocationTopic{},
			storeCarriersLocationTopics:  map[string]*StoreCarriersLocationTopic{},
			globalCarriersLocationTopics: map[string]chan *model.Carrier{},
		},
	}
}

type StoreCarriersTopic struct { //Topicos de store carriers
	Key       string
	Observers map[string]chan *model.Carrier
}
type StoreOrdersTopic struct { //Topicos de store orders
	Key       string
	Observers map[string]chan *model.Order
}
type OrderTopic struct { // Topicos de orders
	Key       string
	Observers map[string]chan *model.Order
}
type CarrierTopic struct { // Topicos de orders
	Key       string
	Observers map[string]chan *model.Carrier
}
type CarrierLocationTopic struct { // Topicos de orders
	Key       string
	Observers map[string]chan *model.Carrier
}
type StoreCarriersLocationTopic struct { // Topicos de orders
	Key       string
	Observers map[string]chan *model.Carrier
}

type Resolver struct {
	sync.Mutex
	carriers                     *mgo.Collection
	orders                       *mgo.Collection
	stores                       *mgo.Collection
	storeCarriersTopics          map[string]*StoreCarriersTopic
	storeOrdersTopics            map[string]*StoreOrdersTopic
	orderTopics                  map[string]*OrderTopic
	carrierTopics                map[string]*CarrierTopic
	carrierLocationTopics        map[string]*CarrierLocationTopic
	storeCarriersLocationTopics  map[string]*StoreCarriersLocationTopic
	globalCarriersLocationTopics map[string]chan *model.Carrier
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
