package mongo

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"errors"
)

func (m *mongo) GetOrganisation(id string, u auth.User) (impact.Organisation, error) {
	org := impact.Organisation{}

	col, closer := m.getOrganisationCollection()
	defer closer()

	userOrg, err := u.Organisation()
	if err != nil {
		return org, err
	}

	if id != userOrg {
		return org, errors.New("User does not have permission to view this organisation")
	}

	err = col.FindId(id).One(&org)
	if err != nil {
		if mgo.ErrNotFound == err {
			return org, data.NewNotFoundError("Organisation")
		}
		return org, err
	}
	return org, nil
}
