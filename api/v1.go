package api

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
)

type meetingTypes struct {
	answerInterface   *graphql.Interface
	intAnswer         *graphql.Object
	categoryAggregate *graphql.Object
	aggregates        *graphql.Object
	meetingType       *graphql.Object
	remoteMeetingType *graphql.Object
}

type organisationTypes struct {
	organisationType *graphql.Object
}

type outcomeSetTypes struct {
	questionInterface *graphql.Interface
	likertScale       *graphql.Object
	outcomeSetType    *graphql.Object
	aggregationEnum   *graphql.Enum
	categoryType      *graphql.Object
}

type reportTypes struct {
	JOCType *graphql.Object
}

type v1 struct {
	db      data.Base
	authGen auth.Generator
}

// NewV1 returns a http.Handler which serves the V1 graphql
func NewV1(db data.Base, authGen auth.Generator) (http.Handler, error) {
	v := &v1{
		db:      db,
		authGen: authGen,
	}
	orgTypes := v.initOrgTypes()
	osTypes := v.initOutcomeSetTypes(orgTypes)
	meetTypes := v.initMeetingTypes(orgTypes, osTypes)
	repTypes := v.initRepTypes()
	schema, err := v.getSchema(orgTypes, osTypes, meetTypes, repTypes)
	if err != nil {
		return nil, err
	}
	h := handler.New(&handler.Config{
		Schema: schema,
		Pretty: true,
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ContextHandler(r.Context(), w, r)
	}), nil
}
