package graph

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/globalsign/mgo"
	"golang.org/x/crypto/bcrypt"
	"github.com/earqq/encargo-backend/graph/model"
	"sync"
    "math/rand"
)
var Observers map[string]chan []*model.Carrier
type Resolver struct {
	sync.Mutex
	carriers *mgo.Collection
	orders   *mgo.Collection
	stores   *mgo.Collection
	observers map[string]chan []*model.Carrier
}
func init(){
	Observers =  map[string]chan []*model.Carrier{}
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