package api

import (
	"errors"
	"github.com/graphql-go/graphql"
	impact "github.com/impactasaurus/server"
)

func (v *v1) initSchemaTypes() {
	v.organisationType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Organisation",
		Description: "An organisation",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Organisation's name",
			},
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Unique identifier for the organisation",
			},
		},
	})

	v.questionType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Question",
		Description: "Question belonging to an outcome set",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Unique ID for the question",
			},
			"question": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The question",
			},
			"type": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Question type",
			},
		},
	})

	v.questionInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "QuestionInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"question": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The question",
			},
			"type": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Question type",
			},
		},
		Description: "A definition of a new outcomeset",
	})

	v.outcomeSetType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "OutcomeSet",
		Description: "A set of questions to determine outcomes",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Unique ID",
			},
			"organisationID": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Organisation's unique ID",
			},
			"organisation": &graphql.Field{
				Type:        graphql.NewNonNull(v.organisationType),
				Description: "The owning organisation of the outcome set",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var organisationID string
					switch t := p.Source.(type) {
					case *impact.OutcomeSet:
						organisationID = t.OrganisationID
					case impact.OutcomeSet:
						organisationID = t.OrganisationID
					}
					if organisationID == "" {
						return nil, errors.New("Expected an OutcomeSet when resolving Organisation")
					}
					return v.db.GetOrganisation(organisationID)
				},
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Name of the outcome set",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "Information about the outcome set",
			},
			"questions": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(v.questionType)),
				Description: "Questions associated with the outcome set",
			},
		},
	})

	v.outcomeSetInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "OutcomeSetInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"organisationID": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The ID of the owning organisation of the outcome set",
			},
			"name": &graphql.InputObjectFieldConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Name of the outcome set",
			},
			"description": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Information about the outcome set",
			},
		},
		Description: "A definition of a new outcomeset",
	})

	v.answerType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Answer",
		Description: "Answer to a question",
		Fields: graphql.Fields{
			"questionID": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The ID of the question answered",
			},
			"answer": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The provided answer",
			},
		},
	})

	v.meetingType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Meeting",
		Description: "A set of answers for an outcome set",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Unique ID for the meeting",
			},
			"beneficiary": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The beneficiary providing the answers",
			},
			"user": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The user who collected the answers",
			},
			"outcomeSetID": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The ID of the outcome set answered",
			},
			"outcomeSet": &graphql.Field{
				Type:        graphql.NewNonNull(v.outcomeSetType),
				Description: "The outcome set answered",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var outcomeSetID string
					switch t := p.Source.(type) {
					case *impact.Meeting:
						outcomeSetID = t.OutcomeSetID
					case impact.Meeting:
						outcomeSetID = t.OutcomeSetID
					}
					if outcomeSetID == "" {
						return nil, errors.New("Expected an Meeting when resolving outcomeSet")
					}
					return v.db.GetOutcomeSet(outcomeSetID)
				},
			},
			"organisationID": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Organisation's unique ID",
			},
			"organisation": &graphql.Field{
				Type:        graphql.NewNonNull(v.organisationType),
				Description: "The owning organisation of the outcome set",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var organisationID string
					switch t := p.Source.(type) {
					case *impact.Meeting:
						organisationID = t.OrganisationID
					case impact.Meeting:
						organisationID = t.OrganisationID
					}
					if organisationID == "" {
						return nil, errors.New("Expected an Meeting when resolving Organisation")
					}
					return v.db.GetOrganisation(organisationID)
				},
			},
			"answers": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(v.answerType)),
				Description: "The answers provided in the meeting",
			},
			"conducted": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "When the meeting was conducted",
			},
			"created": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "When the meeting was entered into the system",
			},
			"modified": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "When the meeting was last modified in the system",
			},
		},
	})
}
