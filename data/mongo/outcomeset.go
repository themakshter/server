package mongo

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) GetOutcomeSet(id string) (*impact.OutcomeSet, error) {
	col, closer := m.getOutcomeCollection()
	defer closer()

	os := &impact.OutcomeSet{}
	err := col.FindId(id).One(os)
	if err != nil {
		if mgo.ErrNotFound == err {
			return nil, &data.NotFound{}
		}
		return nil, err
	}
	return os, nil
}

func (m *mongo) GetOutcomeSets() ([]impact.OutcomeSet, error) {
	col, closer := m.getOutcomeCollection()
	defer closer()

	results := []impact.OutcomeSet{}
	err := col.Find(bson.M{}).All(&results)
	return results, err
}

func (m *mongo) GetQuestion(outcomeSetID string, questionID string) (*impact.Question, error) {
	os, err := m.GetOutcomeSet(outcomeSetID)
	if err != nil {
		return nil, err
	}

	for _, q := range os.Questions {
		if q.ID == questionID {
			return &q, nil
		}
	}
	return nil, &data.NotFound{}
}
