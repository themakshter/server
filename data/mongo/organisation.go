package mongo

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"errors"
)

func (m *mongo) GetOrganisation(id string, u auth.User) (*impact.Organisation, error) {
	col, closer := m.getOrganisationCollection()
	defer closer()

	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}

	if id != userOrg {
		return nil, errors.New("User does not have permission to view this organisation")
	}

	org := &impact.Organisation{}
	err = col.FindId(id).One(org)
	if err != nil {
		if mgo.ErrNotFound == err {
			return nil, &data.NotFound{}
		}
		return nil, err
	}
	return org, nil
}
