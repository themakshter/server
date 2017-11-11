package mongo

import (
	"errors"

	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	uuid "github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) GetOutcomeSet(id string, u auth.User) (impact.OutcomeSet, error) {
	os := impact.OutcomeSet{}

	col, closer := m.getOutcomeCollection()
	defer closer()

	userOrg, err := u.Organisation()
	if err != nil {
		return os, err
	}

	err = col.Find(bson.M{
		"_id":            id,
		"organisationID": userOrg,
	}).One(&os)
	if err != nil {
		if mgo.ErrNotFound == err {
			return os, data.NewNotFoundError("Outcome Set")
		}
		return os, err
	}
	return os, nil
}

func (m *mongo) GetOutcomeSets(u auth.User) ([]impact.OutcomeSet, error) {
	col, closer := m.getOutcomeCollection()
	defer closer()

	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}

	results := []impact.OutcomeSet{}
	err = col.Find(bson.M{
		"organisationID": userOrg,
		"deleted":        false,
	}).All(&results)

	return results, err
}

func (m *mongo) NewOutcomeSet(name, description string, skippable bool, u auth.User) (impact.OutcomeSet, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return impact.OutcomeSet{}, err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	existing, err := col.Find(bson.M{
		"name":           name,
		"organisationID": userOrg,
		"deleted":        false,
	}).Count()
	if err != nil {
		return impact.OutcomeSet{}, err
	}
	if existing != 0 {
		return impact.OutcomeSet{}, errors.New("Name already in use")
	}

	id := uuid.NewV4()

	newOS := impact.OutcomeSet{
		ID:             id.String(),
		Deleted:        false,
		Description:    description,
		Name:           name,
		OrganisationID: userOrg,
		Skippable:      skippable,
	}
	if err := col.Insert(newOS); err != nil {
		return impact.OutcomeSet{}, err
	}
	return m.GetOutcomeSet(id.String(), u)
}

func (m *mongo) EditOutcomeSet(id, name, description string, skippable bool, u auth.User) (impact.OutcomeSet, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return impact.OutcomeSet{}, err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	if err = col.Update(bson.M{
		"_id":            id,
		"organisationID": userOrg,
	}, bson.M{
		"$set": bson.M{
			"name":        name,
			"description": description,
			"skippable":   skippable,
		},
	}); err != nil {
		return impact.OutcomeSet{}, err
	}

	return m.GetOutcomeSet(id, u)
}

func (m *mongo) DeleteOutcomeSet(id string, u auth.User) error {
	userOrg, err := u.Organisation()
	if err != nil {
		return err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	return col.Update(bson.M{
		"_id":            id,
		"organisationID": userOrg,
	}, bson.M{
		"$set": bson.M{
			"deleted": true,
		},
	})
}
