package mongo

import (
	"fmt"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"time"
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
		Key:    []string{"organisationID", "name"},
	}); err != nil {
		return err
	}

	return nil
}
