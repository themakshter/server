package mongo

import (
	"errors"
	"fmt"
	"time"

	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Config is the configuration required to connect to the mongo DB
type Config struct {
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

func New(hostname string, port int, database, user, password string) (data.Base, error) {
	url := fmt.Sprint(hostname, ":", port)
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{url},
		Timeout:  time.Duration(60) * time.Second,
		Database: database,
		Source:   database,
		Username: user,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	m := &mongo{
		baseSession: session,
	}
	if err := m.ensureIndexes(); err != nil {
		return nil, err
	}
	return m, nil
}

type mongo struct {
	baseSession *mgo.Session
}

func (m *mongo) ensureIndexes() error {
	osCol, osCloser := m.getOutcomeCollection()
	defer osCloser()

	if err := osCol.EnsureIndex(mgo.Index{
		Key: []string{"organisationID", "name"},
	}); err != nil {
		return err
	}

	return nil
}

// query adds the addQueryFields method to bson.M
type query bson.M

// addQueryFields ensures that additions to a bson.M does not overwrite previous data
func (q query) addQueryFields(fields map[string]interface{}) error {
	for k, v := range fields {
		if forcedV, set := q[k]; set && forcedV != v {
			return errors.New("Not authorized to make this request")
		}
		q[k] = v
	}
	return nil
}
