package api

import (
	"github.com/graphql-go/graphql"
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
			//"AddQuestion",
			//"EditQuestion",
			//"DeleteQuestion",
			"AddOutcomeSet": &graphql.Field{
				Type:        v.outcomeSetType,
				Description: "Create a new outcomeset",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
						Description: "The name of the outcomeset",
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
						Description: "An optional description",
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					name := p.Args["name"].(string)
					description := ""
					descriptionRaw := p.Args["description"]
					if descriptionRaw != nil {
						description = descriptionRaw.(string)
					}
					return v.db.NewOutcomeSet(name, description, u)
				}),
			},
			"EditOutcomeSet": &graphql.Field{
				Type:        v.outcomeSetType,
				Description: "Edit an outcomeset",
				Args: graphql.FieldConfigArgument{
					"outcomeSetID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
						Description: "The ID of the outcomeset",
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
						Description: "The new name to apply to the outcomeset",
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
						Description: "The new description to apply to the outcomeset, if left null, any existing description will be removed",
					},
				},
				Resolve: userRestrictedResolver(func(p graphql.ResolveParams, u auth.User) (interface{}, error) {
					id := p.Args["outcomeSetID"].(string)
					name := p.Args["name"].(string)
					description := ""
					descriptionRaw := p.Args["description"]
					if descriptionRaw != nil {
						description = descriptionRaw.(string)
					}
					return v.db.EditOutcomeSet(id, name, description, u)
				}),
			},
			//"DeleteOutcomeSet",
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
