package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/earqq/encargo-backend/auth"
	"github.com/earqq/encargo-backend/db"
	"github.com/earqq/encargo-backend/graph"
	"github.com/earqq/encargo-backend/graph/generated"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

const defaultPort = "8096"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error main loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	db.ConnectDB()
	router := chi.NewRouter()
	router.Use(auth.Middleware())
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
	})

	srv := handler.New(generated.NewExecutableSchema(graph.New()))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", c.Handler(srv))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
