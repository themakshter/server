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

	beneficiaryAuthGen := auth.NewBeneficiaryJWTGenerator(c.Local.Audience, c.Local.Issuer, auth.MustParseRSAPrivateKeyFromPEM(c.Local.PrivateKey))
	v1Handler, err := api.NewV1(db, beneficiaryAuthGen)
	if err != nil {
		log.Fatal(err, nil)
	}

	auth0Auth := auth.NewJWTAuthenticator(c.Auth0.Audience, c.Auth0.Issuer, auth.MustParseRSAPublicKeyFromPEM(c.Auth0.PublicKey))
	localAuth := auth.NewBeneficiaryAuthenticator(c.Local.Audience, c.Local.Issuer, auth.MustParseRSAPublicKeyFromPEM(c.Local.PublicKey))
	cors := corsLib.New(corsLib.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	http.Handle("/v1/graphql", cors.Handler(auth.Middleware(v1Handler, auth0Auth, localAuth)))

	if err = http.ListenAndServe(":"+strconv.Itoa(c.Network.Port), nil); err != nil {
		log.Fatal(err, nil)
	}
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
