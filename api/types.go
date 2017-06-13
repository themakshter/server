package api

import (
	"errors"
	"github.com/graphql-go/graphql"
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
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

	v.questionInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name:        "QuestionInterface",
		Description: "The interface satisfied by all question types",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Unique ID for the question",
			},
			"question": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The question",
			},
			"deleted": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Whether the question has been deleted",
			},
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			obj, ok := p.Value.(impact.Question)
			if !ok {
				return v.likertScale
			}
			switch obj.Type {
			case impact.LIKERT:
				return v.likertScale
			default:
				return v.likertScale
			}
		},
	})

	v.likertScale = graphql.NewObject(graphql.ObjectConfig{
		Name:        "LikertScale",
		Description: "Question gathering information using Likert Scales",
		Interfaces: []*graphql.Interface{
			v.questionInterface,
		},
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Unique ID for the question",
			},
			"question": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The question",
			},
			"deleted": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Whether the question has been deleted",
			},
			"minValue": &graphql.Field{
				Type:        graphql.Int,
				Description: "The minimum value in the scale",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					obj, ok := p.Source.(impact.Question)
					if !ok {
						return nil, errors.New("Expecting an impact.Question")
					}
					minValue, ok := obj.Options["minValue"]
					if !ok {
						return nil, nil
					}
					minValueInt, ok := minValue.(int)
					if !ok {
						return nil, errors.New("Min likert value should be an int")
					}
					return minValueInt, nil
				},
			},
			"maxValue": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The maximum value in the scale",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					obj, ok := p.Source.(impact.Question)
					if !ok {
						return nil, errors.New("Expecting an impact.Question")
					}
					maxValue, ok := obj.Options["maxValue"]
					if !ok {
						return nil, nil
					}
					maxValueInt, ok := maxValue.(int)
					if !ok {
						return nil, errors.New("Max likert value should be an int")
					}
					return maxValueInt, nil
				},
			},
			"minLabel": &graphql.Field{
				Type:        graphql.String,
				Description: "The string labelling the minimum value in the scale",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					obj, ok := p.Source.(impact.Question)
					if !ok {
						return nil, errors.New("Expecting an impact.Question")
					}
					label, ok := obj.Options["minLabel"]
					if !ok {
						return nil, nil
					}
					labelStr, ok := label.(string)
					if !ok {
						return nil, errors.New("Min likert label should be an string")
					}
					return labelStr, nil
				},
			},
			"maxLabel": &graphql.Field{
				Type:        graphql.String,
				Description: "The string labelling the maximum value in the scale",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					obj, ok := p.Source.(impact.Question)
					if !ok {
						return nil, errors.New("Expecting an impact.Question")
					}
					label, ok := obj.Options["maxLabel"]
					if !ok {
						return nil, nil
					}
					labelStr, ok := label.(string)
					if !ok {
						return nil, errors.New("Max likert label should be an string")
					}
					return labelStr, nil
				},
			},
		},
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
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					var organisationID string
					switch t := p.Source.(type) {
					case *impact.OutcomeSet:
						organisationID = t.OrganisationID
					case impact.OutcomeSet:
						organisationID = t.OrganisationID
					default:
						return nil, errors.New("Expected an OutcomeSet when resolving Organisation")
					}
					if organisationID == "" {
						return nil, errors.New("Organisation ID is missing on OutcomeSet")
					}
					return v.db.GetOrganisation(organisationID, u)
				}),
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
				Type:        graphql.NewNonNull(graphql.NewList(v.questionInterface)),
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

	v.answerInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name:        "AnswerInterface",
		Description: "The interface satisfied by all answer types",
		Fields: graphql.Fields{
			"questionID": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The ID of the question answered",
			},
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			obj, ok := p.Value.(impact.Answer)
			if !ok {
				return v.numericAnswer
			}
			switch obj.Type {
			case impact.NUMERIC:
				return v.numericAnswer
			default:
				return v.numericAnswer
			}
		},
	})

	v.numericAnswer = graphql.NewObject(graphql.ObjectConfig{
		Name:        "NumericAnswer",
		Description: "Answer containing numeric value",
		Interfaces: []*graphql.Interface{
			v.answerInterface,
		},
		Fields: graphql.Fields{
			"questionID": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The ID of the question answered",
			},
			"answer": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Float),
				Description: "The provided numeric answer",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					obj, ok := p.Source.(impact.Answer)
					if !ok {
						return nil, errors.New("Expecting an impact.Answer")
					}
					num, ok := obj.Answer.(float64)
					if !ok {
						return nil, errors.New("Expected a numeric value")
					}
					return num, nil
				},
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
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					var outcomeSetID string
					switch t := p.Source.(type) {
					case *impact.Meeting:
						outcomeSetID = t.OutcomeSetID
					case impact.Meeting:
						outcomeSetID = t.OutcomeSetID
					default:
						return nil, errors.New("Expected an Meeting when resolving outcomeSet")
					}
					if outcomeSetID == "" {
						return nil, errors.New("OutcomeSetID is missing on Meeting")
					}
					return v.db.GetOutcomeSet(outcomeSetID, u)
				}),
			},
			"organisationID": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Organisation's unique ID",
			},
			"organisation": &graphql.Field{
				Type:        graphql.NewNonNull(v.organisationType),
				Description: "The owning organisation of the outcome set",
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					var organisationID string
					switch t := p.Source.(type) {
					case *impact.Meeting:
						organisationID = t.OrganisationID
					case impact.Meeting:
						organisationID = t.OrganisationID
					default:
						return nil, errors.New("Expected an Meeting when resolving Organisation")
					}
					if organisationID == "" {
						return nil, errors.New("OrganisationID is missing on Meeting")
					}
					return v.db.GetOrganisation(organisationID, u)
				}),
			},
			"answers": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(v.answerInterface)),
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
