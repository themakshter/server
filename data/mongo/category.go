package mongo

import (
	"errors"
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) GetCategory(outcomeSetID, categoryID string, u auth.User) (impact.Category, error) {
	os, err := m.GetOutcomeSet(outcomeSetID, u)
	if err != nil {
		return impact.Category{}, err
	}

	for _, c := range os.Categories {
		if c.ID == categoryID {
			return c, nil
		}
	}
	return impact.Category{}, data.NewNotFoundError("Category")
}

func (m *mongo) NewCategory(outcomeSetID, name, description, aggregation string, u auth.User) (impact.Category, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return impact.Category{}, err
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	id := uuid.NewV4()

	newCategory := &impact.Category{
		ID:          id.String(),
		Name:        name,
		Description: description,
		Aggregation: aggregation,
	}

	if err := col.Update(bson.M{
		"_id":            outcomeSetID,
		"organisationID": userOrg,
	}, bson.M{
		"$push": bson.M{
			"categories": newCategory,
		},
	}); err != nil {
		return impact.Category{}, err
	}

	return m.GetCategory(outcomeSetID, id.String(), u)
}

func (m *mongo) DeleteCategory(outcomeSetID, categoryID string, u auth.User) error {
	userOrg, err := u.Organisation()
	if err != nil {
		return err
	}

	os, err := m.GetOutcomeSet(outcomeSetID, u)
	if err != nil {
		return err
	}

	catQuestions := os.GetCategoryQuestions(categoryID)
	if len(catQuestions) > 0 {
		return errors.New("Cannot delete a category which is being used")
	}

	col, closer := m.getOutcomeCollection()
	defer closer()

	return col.Update(bson.M{
		"_id":            outcomeSetID,
		"organisationID": userOrg,
	}, bson.M{
		"$pull": bson.M{
			"categories": bson.M{
				"id": categoryID,
			},
		},
	})
}
