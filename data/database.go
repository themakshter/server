package data

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
)

// NotFound is an error struct to be used when request to fetch / update / delete an individual element fails because the element does not exist.
type NotFound struct {
}

func (*NotFound) Error() string {
	return "Not Found"
}

type Base interface {
	GetOutcomeSet(id string, u auth.User) (*impact.OutcomeSet, error)
	GetOutcomeSets(u auth.User) ([]impact.OutcomeSet, error)
	GetQuestion(outcomeSetID string, questionID string, u auth.User) (*impact.Question, error)

	GetOrganisation(id string, u auth.User) (*impact.Organisation, error)

	GetMeeting(id string, u auth.User) (*impact.Meeting, error)
	GetMeetingsForBeneficiary(beneficiary string, u auth.User) ([]impact.Meeting, error)
}
