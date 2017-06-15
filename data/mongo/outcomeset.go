package mongo

import (
	"errors"
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func filterOutDeletedQuestions(os *impact.OutcomeSet) *impact.OutcomeSet {
	newQs := make([]impact.Question, 0, len(os.Questions))
	for _, q := range os.Questions {
		if !q.Deleted {
			newQs = append(newQs, q)
		}
	}
	os.Questions = newQs
	return os
}

func (m *mongo) GetOutcomeSet(id string, u auth.User) (*impact.OutcomeSet, error) {
	col, closer := m.getOutcomeCollection()
	defer closer()

	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}

	os := &impact.OutcomeSet{}
	err = col.Find(bson.M{
		"_id":            id,
		"organisationID": userOrg,
	}).One(os)
	if err != nil {
		if mgo.ErrNotFound == err {
			return nil, data.NewNotFoundError("Outcome Set")
		}
		return nil, err
	}
	return filterOutDeletedQuestions(os), nil
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

	for idx, os := range results {
		results[idx] = *filterOutDeletedQuestions(&os)
	}

	return results, err
}

func (m *mongo) NewOutcomeSet(name, description string, u auth.User) (*impact.OutcomeSet, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	existing, err := col.Find(bson.M{
		"name":           name,
		"organisationID": userOrg,
		"deleted":        false,
	}).Count()
	if err != nil {
		return nil, err
	}
	if existing != 0 {
		return nil, errors.New("Name already in use")
	}

	id := uuid.NewV4()

	newOS := &impact.OutcomeSet{
		ID:             id.String(),
		Deleted:        false,
		Description:    description,
		Name:           name,
		OrganisationID: userOrg,
	}
	if err := col.Insert(newOS); err != nil {
		return nil, err
	}
	return m.GetOutcomeSet(id.String(), u)
}

func (m *mongo) EditOutcomeSet(id, name, description string, u auth.User) (*impact.OutcomeSet, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
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
		},
	}); err != nil {
		return nil, err
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
