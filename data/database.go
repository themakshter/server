package data

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"fmt"
)

// NotFound is an error struct to be used when request to fetch / update / delete an individual element fails because the element does not exist.
type notFound struct {
	thing string
}

func NewNotFoundError(thing string) error {
	return &notFound{thing}
}

func (nf *notFound) Error() string {
	return fmt.Sprintf("%s not Found", nf.thing)
}

type Base interface {
	NewOutcomeSet(name, description string, u auth.User) (*impact.OutcomeSet, error)
	GetOutcomeSet(id string, u auth.User) (*impact.OutcomeSet, error)
	GetOutcomeSets(u auth.User) ([]impact.OutcomeSet, error)
	GetQuestion(outcomeSetID string, questionID string, u auth.User) (*impact.Question, error)

	GetOrganisation(id string, u auth.User) (*impact.Organisation, error)

	GetMeeting(id string, u auth.User) (*impact.Meeting, error)
	GetMeetingsForBeneficiary(beneficiary string, u auth.User) ([]impact.Meeting, error)
}
