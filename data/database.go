package data

import impact "github.com/impactasaurus/server"

type Base interface {
	GetOutcomeSet(id string) (*impact.OutcomeSet, error)
	GetOutcomeSets() ([]impact.OutcomeSet, error)
	GetQuestion(outcomeSetID string, questionID string) (*impact.Question, error)

	GetOrganisation(id string) (*impact.Organisation, error)

	GetMeeting(id string) (*impact.Meeting, error)
	GetMeetingsForBeneficiary(beneficiary string) ([]impact.Meeting, error)
}
