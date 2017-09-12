package main

import (
	"net/http"

	"strconv"

	"github.com/impactasaurus/server/api"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data/mongo"
	"github.com/impactasaurus/server/log"
	corsLib "github.com/rs/cors"
)

func main() {
	c := mustGetConfiguration()

	mustConfigureLogger(c)

	db, err := mongo.New(c.Mongo.URL, c.Mongo.Port, c.Mongo.Database, c.Mongo.User, c.Mongo.Password)
	if err != nil {
		log.Fatal(err, nil)
	}

	v1Handler, err := api.NewV1(db)
	if err != nil {
		log.Fatal(err, nil)
	}

	cors := corsLib.New(corsLib.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	http.Handle("/v1/graphql", cors.Handler(auth.Middleware(v1Handler)))

	http.ListenAndServe(":"+strconv.Itoa(c.Network.Port), nil)
}

func mustConfigureLogger(c *config) {
	if c.Sentry.DSN != "" {
		s, err := log.NewSentryErrorTracker(c.Sentry.DSN)
		if err != nil {
			log.Fatal(err, nil)
		}
		log.RegisterErrorTracker(s)
	}
}
