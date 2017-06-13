package api

import (
	"github.com/graphql-go/graphql"
	"github.com/impactasaurus/server/auth"
)

type userAuthenticatedResolver func(graphql.ResolveParams, auth.User) (interface{}, error)

func userRestrictedResolver(fn userAuthenticatedResolver) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		u, e := auth.GetUser(p.Context)
		if e != nil {
			return nil, e
		}
		return fn(p, u)
	}
}