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
			carriers: db.GetCollection("carriers"),
			orders:   db.GetCollection("orders"),
			stores:   db.GetCollection("stores"),
		},
	}
}

var Observers map[string]chan []*model.Carrier
var CarrierOrdersObserver map[string]chan *model.Order
var StoreOrdersObserver map[string]chan *model.Order

type Resolver struct {
	sync.Mutex
	carriers  *mgo.Collection
	orders    *mgo.Collection
	stores    *mgo.Collection
	observers map[string]chan []*model.Carrier
}

func init() {
	Observers = map[string]chan []*model.Carrier{}
	CarrierOrdersObserver = map[string]chan *model.Order{}
	StoreOrdersObserver = map[string]chan *model.Order{}
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
