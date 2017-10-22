package main

import (
	"net/http"

	"strconv"

	jwtLib "github.com/dgrijalva/jwt-go"
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

	auth0Auth := mustGetAuthenticator(c.Auth0)
	cors := corsLib.New(corsLib.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	http.Handle("/v1/graphql", cors.Handler(auth.Middleware(v1Handler, auth0Auth)))

	http.ListenAndServe(":"+strconv.Itoa(c.Network.Port), nil)
}

func mustGetAuthenticator(c configAuth) auth.Authenticator {
	pubKey, err := jwtLib.ParseRSAPublicKeyFromPEM([]byte(c.PublicKey))
	if err != nil {
		log.Fatal(err, nil)
	}
	return auth.NewJWTAuthenticator(c.Audience, c.Issuer, pubKey)
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
