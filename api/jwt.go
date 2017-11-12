package api

import (
	"github.com/graphql-go/graphql"
)

func (v *v1) getJWTQueries() graphql.Fields {
	return graphql.Fields{
		"jwt": &graphql.Field{
			Type:        graphql.String,
			Description: "Get a JWT given a JTI",
			Args: graphql.FieldConfigArgument{
				"jti": &graphql.ArgumentConfig{
					Description: "The ID of the JWT",
					Type:        graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return v.db.GetJWT(p.Args["jti"].(string))
			},
		},
	}
}
