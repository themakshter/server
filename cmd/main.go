package main

import (
	"net/http"

	"github.com/impactasaurus/server/api"
	"github.com/impactasaurus/server/data"
	"github.com/rs/cors"
)

func main() {
	db := data.NewMongo()

	v1Handler, err := api.NewV1(db)
	if err != nil {
		panic(err)
	}

	c := cors.New(cors.Options{
		AllowCredentials: true,
	})
	http.Handle("/v1/graphql", c.Handler(v1Handler))

	http.ListenAndServe(":80", nil)
}
