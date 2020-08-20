package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/earqq/encargo-backend/db"
	"github.com/earqq/encargo-backend/graph/model"
	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2/bson"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	Name string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file auth")
	}
	privateBytes, err := ioutil.ReadFile(os.Getenv("PROJECT_PATH") + "/private.rsa")
	if err != nil {
		log.Fatal("No se puede leer llave privada")
	}
	publicBytes, err := ioutil.ReadFile(os.Getenv("PROJECT_PATH") + "/public.rsa.pub")
	if err != nil {
		log.Fatal("No se puedo leer llave pública")
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se pudo parsear llave privada")
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("No se pudo parsear llave pública")
	}
}

// Generar un nuevo token
func GenerateJWT(username string, userType string) string {
	claims := model.Claim{
		Username: username,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			Issuer: "Login token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, _ := token.SignedString(privateKey)
	return result
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			if err == nil && token.Valid {
				var carrier model.Carrier
				var store model.Store
				tokenString := TokenFromHttpRequest(r)
				usernameFromToken := UsernameFromToken(tokenString)
				userTypeFromToken := UserTypeFromToken(tokenString)
				// Loguear al usuario dependiendo si es usuario tienda o de repartidor
				ctx := context.WithValue(r.Context(), userCtxKey, &carrier)
				if userTypeFromToken == "carrier" {
					carriersDB := db.GetCollection("carriers")
					_ = carriersDB.Find(bson.M{"username": usernameFromToken}).Select(bson.M{"password": 0}).One(&carrier)
					var user = model.User{
						carrier.ID,
						"carrier",
						usernameFromToken,
					}
					ctx = context.WithValue(r.Context(), userCtxKey, &user)
				} else {
					storeDB := db.GetCollection("stores")
					_ = storeDB.Find(bson.M{"username": usernameFromToken}).Select(bson.M{"password": 0}).One(&store)
					var user = model.User{
						store.ID,
						"store",
						usernameFromToken,
					}
					ctx = context.WithValue(r.Context(), userCtxKey, &user)
				}
				// put it in context
				// and call the next with our new context
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)

		})
	}
}

// Obtener token en formato string desde el header
func TokenFromHttpRequest(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	var tokenString string
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}
	return tokenString
}
func JwtDecode(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}

// Obtener tipo de usuario
func UserTypeFromToken(tokenString string) string {
	token, err := JwtDecode(tokenString)
	if err != nil {
		fmt.Println(err)
		return "1"
	}
	if claims, ok := token.Claims.(*model.Claim); ok && token.Valid {
		if claims == nil {
			return "2 "
		}
		return claims.UserType
	} else {
		return "3"
	}
}

//Obtener username de usuario para login
func UsernameFromToken(tokenString string) string {

	token, err := JwtDecode(tokenString)
	if err != nil {
		fmt.Println(err)
		return "1"
	}
	if claims, ok := token.Claims.(*model.Claim); ok && token.Valid {
		if claims == nil {
			return "2 "
		}
		return claims.Username
	} else {
		return "3"
	}
}
func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}

// Obtener usuario logueado
func GetAuthFromContext(ctx context.Context) *model.User {
	return ForContext(ctx)
}
