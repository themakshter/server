package data

import impact "github.com/impactasaurus/server"

// NotFound is an error struct to be used when request to fetch / update / delete an individual element fails because the element does not exist.
type NotFound struct {
}

func (*NotFound) Error() string {
	return "Not Found"
}

type Base interface {
	GetOutcomeSet(id string) (*impact.OutcomeSet, error)
	GetOutcomeSets() ([]impact.OutcomeSet, error)
	GetQuestion(outcomeSetID string, questionID string) (*impact.Question, error)

	GetOrganisation(id string) (*impact.Organisation, error)

	GetMeeting(id string) (*impact.Meeting, error)
	GetMeetingsForBeneficiary(beneficiary string) ([]impact.Meeting, error)
}
