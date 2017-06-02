package main

import "github.com/kelseyhightower/envconfig"

type configMongo struct {
	// User is the username of the mongodb user, leave blank if username and password is not required
	User string `envconfig:"MONGO_USER"`
	// Password is the mongodb user's password
	Password string `envconfig:"MONGO_PASS"`
	// URL is the mongo database's URL
	URL string `envconfig:"MONGO_URL" required:"true"`
	// Port is the network port on which the mongo database is listening
	Port int `envconfig:"MONGO_PORT" required:"true"`
	// Database is the name of the mongo database to use
	Database string `envconfig:"MONGO_DB" required:"true"`
}

type config struct {
	Mongo configMongo
}

func mustGetConfiguration() *config {
	c := &config{}
	envconfig.MustProcess("", c)
	return c
}
