package server

import (
	"time"
)

type Answer struct {
	QuestionID string      `json:"questionID" bson:"questionID"`
	Answer     interface{} `json:"answer"`
}

type Meeting struct {
	ID             string    `json:"id" bson:"_id"`
	Beneficiary    string    `json:"beneficiary"`
	User           string    `json:"user"`
	OutcomeSetID   string    `json:"outcomeSetID" bson:"outcomeSetID"`
	OrganisationID string    `json:"organisationID" bson:"organisationID"`
	Answers        []Answer  `json:"answers"`
	Conducted      time.Time `json:"conducted"`
	Created        time.Time `json:"created"`
	Modified       time.Time `json:"modified"`
}
