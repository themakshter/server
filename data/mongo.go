package data

import impact "github.com/impactasaurus/server"

type mongo struct {
}

func NewMongo() Base {
	return &mongo{}
}

func (m *mongo) GetOutcomeSet(id string) (*impact.OutcomeSet, error) {
	return &impact.OutcomeSet{
		ID:             id,
		OrganisationID: "test",
		Name:           "dummy just built with ID",
	}, nil
}
func (m *mongo) GetOutcomeSets() ([]impact.OutcomeSet, error) {
	return []impact.OutcomeSet{{
		ID:             "2",
		OrganisationID: "test2",
		Name:           "testing 2",
	}, {
		ID:             "22",
		OrganisationID: "test22",
		Name:           "testing 22",
	}}, nil
}
func (m *mongo) GetQuestion(outcomeSetID string, questionID string) (*impact.Question, error) {
	return &impact.Question{
		ID:       questionID,
		Question: "dummy with the given ID",
		Type:     "scale",
	}, nil
}

func (m *mongo) GetOrganisation(id string) (*impact.Organisation, error) {
	return &impact.Organisation{
		ID:   id,
		Name: "test with given ID",
	}, nil
}

func (m *mongo) GetMeeting(id string) (*impact.Meeting, error) {
	return &impact.Meeting{
		ID:             id,
		Beneficiary:    "test with given ID",
		OrganisationID: "test",
		OutcomeSetID:   "out22",
		Answers: []impact.Answer{
			impact.Answer{
				QuestionID: "guid-q-1",
				Answer:     "test answer",
			},
		},
	}, nil
}

func (m *mongo) GetMeetingsForBeneficiary(beneficiary string) ([]impact.Meeting, error) {
	return []impact.Meeting{
		impact.Meeting{
			ID:             "guid-m-1",
			Beneficiary:    beneficiary,
			OrganisationID: "test",
			OutcomeSetID:   "out25",
			Answers: []impact.Answer{
				impact.Answer{
					QuestionID: "guid-q-1",
					Answer:     "test answer",
				},
			},
		},
	}, nil
}
