package api

import (
	"github.com/graphql-go/graphql"
)

func (v *v1) getSchema() (*graphql.Schema, error) {

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"outcomesets": &graphql.Field{
				Type: graphql.NewList(v.outcomeSetType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return v.db.GetOutcomeSets()
				},
			},
			"outcomeset": &graphql.Field{
				Type: v.outcomeSetType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "The ID of the outcomeset",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return v.db.GetOutcomeSet(p.Args["id"].(string))
				},
			},
			"organisation": &graphql.Field{
				Type: v.organisationType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "The ID of the organisation",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return v.db.GetOrganisation(p.Args["id"].(string))
				},
			},
			"meeting": &graphql.Field{
				Type: v.meetingType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "The ID of the meeting",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return v.db.GetMeeting(p.Args["id"].(string))
				},
			},
			"meetings": &graphql.Field{
				Type: graphql.NewList(v.meetingType),
				Args: graphql.FieldConfigArgument{
					"beneficiary": &graphql.ArgumentConfig{
						Description: "The ID of the beneficiary",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return v.db.GetMeetingsForBeneficiary(p.Args["beneficiary"].(string))
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"AddOutcomeSet": &graphql.Field{
				Type:        v.outcomeSetType,
				Description: "Create an outcomeset",
				Args: graphql.FieldConfigArgument{
					"outcomesetIn": &graphql.ArgumentConfig{
						Type:        v.outcomeSetInputType,
						Description: "The new outcomeset",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					os := p.Args["outcomesetIn"].(map[string]interface{})
					return v.db.GetOutcomeSet(os["organisation"].(string))
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
		Types: []graphql.Type{
			v.likertScale,
		},
	})
	if err != nil {
		return nil, err
	}
	return &schema, nil
}
