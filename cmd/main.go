package main

import (
	"net/http"

	"github.com/impactasaurus/server/api"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data/mongo"
	corsLib "github.com/rs/cors"
)

func main() {
	c := mustGetConfiguration()

	db, err := mongo.New(c.Mongo.URL, c.Mongo.Port, c.Mongo.Database, c.Mongo.User, c.Mongo.Password)
	if err != nil {
		panic(err)
	}

	v1Handler, err := api.NewV1(db)
	if err != nil {
		panic(err)
	}

	cors := corsLib.New(corsLib.Options{
		AllowCredentials: true,
	})
	http.Handle("/v1/graphql", cors.Handler(auth.Middleware(v1Handler)))

	http.ListenAndServe(":80", nil)
}
