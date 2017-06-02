package mongo

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
)

func (m *mongo) GetOrganisation(id string) (*impact.Organisation, error) {
	col, closer := m.getOrganisationCollection()
	defer closer()

	org := &impact.Organisation{}
	err := col.FindId(id).One(org)
	if err != nil {
		if mgo.ErrNotFound == err {
			return nil, &data.NotFound{}
		}
		return nil, err
	}
	return org, nil
}
