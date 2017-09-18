package mongo

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) GetQuestion(outcomeSetID string, questionID string, u auth.User) (impact.Question, error) {
	os, err := m.GetOutcomeSet(outcomeSetID, u)
	if err != nil {
		return impact.Question{}, err
	}

	for _, q := range os.Questions {
		if q.ID == questionID {
			return q, nil
		}
	}
	return impact.Question{}, data.NewNotFoundError("Question")
}

func (m *mongo) NewQuestion(outcomeSetID, question, description string, questionType impact.QuestionType, options map[string]interface{}, u auth.User) (impact.Question, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return impact.Question{}, err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	id := uuid.NewV4()

	newQuestion := &impact.Question{
		ID:          id.String(),
		Question:    question,
		Description: description,
		Type:        questionType,
		Options:     options,
		Deleted:     false,
	}

	if err := col.Update(bson.M{
		"_id":            outcomeSetID,
		"organisationID": userOrg,
	}, bson.M{
		"$push": bson.M{
			"questions": newQuestion,
		},
	}); err != nil {
		return impact.Question{}, err
	}

	return m.GetQuestion(outcomeSetID, id.String(), u)
}

func (m *mongo) DeleteQuestion(outcomeSetID, questionID string, u auth.User) error {
	userOrg, err := u.Organisation()
	if err != nil {
		return err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	return col.Update(bson.M{
		"_id":            outcomeSetID,
		"organisationID": userOrg,
		"questions.id":   questionID,
	}, bson.M{
		"$set": bson.M{
			"questions.$.deleted": true,
		},
	})
}

func (m *mongo) EditQuestion(outcomeSetID, questionID, question, description string, questionType impact.QuestionType, options map[string]interface{}, u auth.User) (impact.Question, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return impact.Question{}, err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	newQ := impact.Question{
		ID:          questionID,
		Question:    question,
		Description: description,
		Type:        questionType,
		Options:     options,
		Deleted:     false,
	}

	if err := col.Update(bson.M{
		"_id":            outcomeSetID,
		"organisationID": userOrg,
		"questions.id":   questionID,
	}, bson.M{
		"$set": bson.M{
			"questions.$": newQ,
		},
	}); err != nil {
		return impact.Question{}, err
	}
	return newQ, nil
}

func (m *mongo) SetCategory(outcomeSetID, questionID, categoryID string, u auth.User) (impact.Question, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return impact.Question{}, err
	}

	_, err = m.GetCategory(outcomeSetID, categoryID, u)
	if err != nil {
		return impact.Question{}, data.NewNotFoundError("Category")
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	if err := col.Update(bson.M{
		"_id":            outcomeSetID,
		"organisationID": userOrg,
		"questions.id":   questionID,
	}, bson.M{
		"$set": bson.M{
			"questions.$.categoryID": categoryID,
		},
	}); err != nil {
		return impact.Question{}, err
	}
	return m.GetQuestion(outcomeSetID, questionID, u)
}

func (m *mongo) RemoveCategory(outcomeSetID, questionID string, u auth.User) (impact.Question, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return impact.Question{}, err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	if err := col.Update(bson.M{
		"_id":            outcomeSetID,
		"organisationID": userOrg,
		"questions.id":   questionID,
	}, bson.M{
		"$set": bson.M{
			"questions.$.categoryID": nil,
		},
	}); err != nil {
		return impact.Question{}, err
	}
	return m.GetQuestion(outcomeSetID, questionID, u)
}
