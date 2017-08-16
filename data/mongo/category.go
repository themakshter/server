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

func (m *mongo) NewCategory(outcomeSetID, name, description string, aggregation impact.Aggregation, u auth.User) (impact.Category, error) {
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

func (m *mongo) isCategoryActivelyUsed(os impact.OutcomeSet, categoryID string) error {
	unarchivedQuestions := os.GetCategoryQuestions(categoryID)
	if len(unarchivedQuestions) > 0 {
		return errors.New("Cannot delete a category which is being used")
	}
	return nil
}

func (m *mongo) removeCategoryFromArchivedCategoryQuestions(os impact.OutcomeSet, categoryID string, u auth.User) {
	archivedQuestions := os.GetArchivedCategoryQuestions(categoryID)
	for _, q := range archivedQuestions {
		m.RemoveCategory(os.ID, q.ID, u)
	}
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

	if err := m.isCategoryActivelyUsed(os, categoryID); err != nil {
		return err
	}

	m.removeCategoryFromArchivedCategoryQuestions(os, categoryID, u)

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
