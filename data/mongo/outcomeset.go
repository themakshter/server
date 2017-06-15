package mongo

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	uuid "github.com/satori/go.uuid"
)

func (m *mongo) GetOutcomeSet(id string, u auth.User) (*impact.OutcomeSet, error) {
	col, closer := m.getOutcomeCollection()
	defer closer()

	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}

	os := &impact.OutcomeSet{}
	err = col.Find(bson.M{
		"_id": id,
		"organisationID": userOrg,
	}).One(os)
	if err != nil {
		if mgo.ErrNotFound == err {
			return nil, data.NewNotFoundError("Outcome Set")
		}
		return nil, err
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
	}).All(&results)
	return results, err
}

func (m *mongo) GetQuestion(outcomeSetID string, questionID string, u auth.User) (*impact.Question, error) {
	os, err := m.GetOutcomeSet(outcomeSetID, u)
	if err != nil {
		return nil, err
	}

	for _, q := range os.Questions {
		if q.ID == questionID {
			return &q, nil
		}
	}
	return nil, data.NewNotFoundError("Question")
}

func (m *mongo) NewOutcomeSet(name, description string, u auth.User) (*impact.OutcomeSet, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}

	id := uuid.NewV4()

	newOS := &impact.OutcomeSet{
		ID: id.String(),
		Deleted: false,
		Description: description,
		Name: name,
		OrganisationID: userOrg,
	}

	col, closer := m.getOutcomeCollection()
	defer closer()
	if err := col.Insert(newOS); err != nil {
		return nil, err
	}
	return m.GetOutcomeSet(id.String(), u)

}
