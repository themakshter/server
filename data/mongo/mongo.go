package mongo

import (
	"errors"
	"fmt"
	"time"

	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongo struct {
	baseSession *mgo.Session
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
		q["k"] = v
	}
	return nil
}
