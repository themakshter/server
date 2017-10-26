package main

import (
	"github.com/kelseyhightower/envconfig"
	"strings"
	"github.com/impactasaurus/server/data/mongo"
)

type configAuth struct {
	Audience  string `required:"true"`
	Issuer    string `required:"true"`
	PublicKey string `required:"true"`
}

type configAuthGen struct {
	configAuth
	PrivateKey string `required:"true"`
}

type configNetwork struct {
	Port int `envconfig:"PORT" default:"80"`
}

type configErrorTracking struct {
	DSN string `envconfig:"SENTRY_DSN" default:""`
}

type config struct {
	Mongo   mongo.Config
	Network configNetwork
	Sentry  configErrorTracking
	Auth0   configAuth
	Local   configAuthGen
}

func mustGetConfiguration() *config {
	c := &config{}
	envconfig.MustProcess("", c)
	// have to process separately as embedded struct
	envconfig.MustProcess("LOCAL", &c.Local.configAuth)
	// required to correctly parse the new lines in the keys
	c.Auth0.PublicKey = strings.Replace(c.Auth0.PublicKey, "\\n", "\n", -1)
	c.Local.PublicKey = strings.Replace(c.Local.PublicKey, "\\n", "\n", -1)
	c.Local.PrivateKey = strings.Replace(c.Local.PrivateKey, "\\n", "\n", -1)
	return c
}
