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

const aud = "pfKiAOUJh5r6jCxRn5vUYq7odQsjPUKf"
const iss = "https://impact.eu.auth0.com/"
const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA72wlFrwLpR2RV5+TxcQU
iLqAqILTRLLWkYJnw+xePiuSIutXARkOcyKGfsunZ/1xZLe5rRuP06M/RTKCNBrq
WwcSLEjL7Dh4JbrkDCGWx7YFcKzswZ3/3Gb2BtavpDhfg5RzaIBgnC+uyKAZvocA
+YoKT3SPHm79vnmvQvpheMqfuNJSmw0mCXMcwqjJUQ2GLQsWfI2qeVybxcsRDqp5
5kDGPVThrg8OGwJQDrHE4DopXWHWvUxKVS/e6eFN6qgibVP7vD3SnL0M7wgJDQuk
TzskiF5Zzsgc86b2P6kQ1H2ryKp1jNDjkBCpr3F7KfNek/ADrSZVpxFSg4cve7X3
+wIDAQAB
-----END PUBLIC KEY-----`

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

	jwtAuthenticator := auth.NewJWTAuthenticator(aud, iss, publicKey)
	cors := corsLib.New(corsLib.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	http.Handle("/v1/graphql", cors.Handler(auth.Middleware(v1Handler, jwtAuthenticator)))

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
