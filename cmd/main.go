package main

import (
	"net/http"

	"github.com/impactasaurus/server/api"
	"github.com/impactasaurus/server/data"
)

func main() {
	db := data.NewMongo()
	v1 := api.NewV1(db)
	if err := v1.Listen(); err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", nil)
}
