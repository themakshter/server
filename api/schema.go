package api

import (
	"github.com/graphql-go/graphql"
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
)

func (v *v1) getSchema() (*graphql.Schema, error) {

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"outcomesets": &graphql.Field{
				Type:        graphql.NewList(v.outcomeSetType),
				Description: "Gather all outcome sets",
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					return v.db.GetOutcomeSets(u)
				}),
			},
			"outcomeset": &graphql.Field{
				Type:        v.outcomeSetType,
				Description: "Gather a specific outcome set",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "The ID of the outcomeset",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					return v.db.GetOutcomeSet(p.Args["id"].(string), u)
				}),
			},
			"organisation": &graphql.Field{
				Type:        v.organisationType,
				Description: "Get an organisation by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "The ID of the organisation",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					return v.db.GetOrganisation(p.Args["id"].(string), u)
				}),
			},
			"meeting": &graphql.Field{
				Type:        v.meetingType,
				Description: "Get a meeting by meeting ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "The ID of the meeting",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					return v.db.GetMeeting(p.Args["id"].(string), u)
				}),
			},
			"meetings": &graphql.Field{
				Type:        graphql.NewList(v.meetingType),
				Description: "Get all meetings associated with a beneficiary",
				Args: graphql.FieldConfigArgument{
					"beneficiary": &graphql.ArgumentConfig{
						Description: "The ID of the beneficiary",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					return v.db.GetMeetingsForBeneficiary(p.Args["beneficiary"].(string), u)
				}),
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"AddOutcomeSet": &graphql.Field{
				Type:        v.outcomeSetType,
				Description: "Create a new outcomeset",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The name of the outcomeset",
					},
					"description": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "An optional description",
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					name := p.Args["name"].(string)
					description := getNullableString(p.Args, "description")
					return v.db.NewOutcomeSet(name, description, u)
				}),
			},
			"EditOutcomeSet": &graphql.Field{
				Type:        v.outcomeSetType,
				Description: "Edit an outcomeset",
				Args: graphql.FieldConfigArgument{
					"outcomeSetID": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.ID),
						Description: "The ID of the outcomeset",
					},
					"name": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The new name to apply to the outcomeset",
					},
					"description": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "The new description to apply to the outcomeset, if left null, any existing description will be removed",
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					id := p.Args["outcomeSetID"].(string)
					name := p.Args["name"].(string)
					description := getNullableString(p.Args, "description")
					return v.db.EditOutcomeSet(id, name, description, u)
				}),
			},
			"DeleteOutcomeSet": &graphql.Field{
				Type:        graphql.ID,
				Description: "Deletes an outcomeset and returns the ID of the deleted outcomeset",
				Args: graphql.FieldConfigArgument{
					"outcomeSetID": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.ID),
						Description: "The ID of the outcomeset",
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					id := p.Args["outcomeSetID"].(string)
					if err := v.db.DeleteOutcomeSet(id, u); err != nil {
						return nil, err
					}
					return id, nil
				}),
			},
			"AddLikertQuestion": &graphql.Field{
				Type:        v.outcomeSetType,
				Description: "Add a likert scale question to an outcome set",
				Args: graphql.FieldConfigArgument{
					"outcomeSetID": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.ID),
						Description: "The ID of the outcomeset",
					},
					"question": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Question to be asked",
					},
					"minValue": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "Minimum value of the likert scale",
					},
					"maxValue": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: "Maximum value of the likert scale",
					},
					"minLabel": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Label associated with the minimum value of the likert scale",
					},
					"maxLabel": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Label associated with the maximum value of the likert scale",
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					id := p.Args["outcomeSetID"].(string)
					question := p.Args["question"].(string)
					minValue := getNullableInt(p.Args, "minValue")
					maxValue := p.Args["maxValue"].(int)
					minLabel := getNullableString(p.Args, "minLabel")
					maxLabel := getNullableString(p.Args, "maxLabel")
					if _, err := v.db.NewQuestion(id, question, impact.LIKERT, map[string]interface{}{
						"minValue": minValue,
						"maxValue": maxValue,
						"minLabel": minLabel,
						"maxLabel": maxLabel,
					}, u); err != nil {
						return nil, err
					}
					return v.db.GetOutcomeSet(id, u)
				}),
			},
			"DeleteQuestion": &graphql.Field{
				Type:        v.outcomeSetType,
				Description: "Remove a question from an outcome set",
				Args: graphql.FieldConfigArgument{
					"outcomeSetID": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The ID of the outcomeset",
					},
					"questionID": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The ID of the question",
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					outcomeSetID := p.Args["outcomeSetID"].(string)
					questionID := p.Args["questionID"].(string)
					if err := v.db.DeleteQuestion(outcomeSetID, questionID, u); err != nil {
						return nil, err
					}
					return v.db.GetOutcomeSet(outcomeSetID, u)
				}),
			},
			//"AddMeeting",
			//"EditMeeting",
			//"DeleteMeeting",
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
		Types: []graphql.Type{
			v.likertScale,
			v.numericAnswer,
		},
	})
	if err != nil {
		return nil, err
	}
	return &schema, nil
}
