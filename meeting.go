package server

import (
	"time"
)

type Answer struct {
	QuestionID string `json:"questionID"`
	Answer     string `json:"answer"`
}

type Meeting struct {
	ID             string    `json:"id"`
	Beneficiary    string    `json:"beneficiary"`
	User           string    `json:"user"`
	OutcomeSetID   string    `json:"outcomeSetID"`
	OrganisationID string    `json:"organisationID"`
	Answers        []Answer  `json:"answers"`
	Conducted      time.Time `json:"conducted"`
	Created        time.Time `json:"created"`
	Modified       time.Time `json:"modified"`
}
