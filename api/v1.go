package api

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/impactasaurus/server/data"
	"net/http"
)

type v1 struct {
	db data.Base

	questionInterface   *graphql.Interface
	likertScale         *graphql.Object
	outcomeSetType      *graphql.Object
	outcomeSetInputType *graphql.InputObject

	answerInterface *graphql.Interface
	numericAnswer   *graphql.Object
	meetingType     *graphql.Object

	organisationType *graphql.Object
}

type graphqlHandlerWrapper struct {
	h *handler.Handler
}

func (gw *graphqlHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gw.h.ContextHandler(r.Context(), w, r)
}

func NewV1(db data.Base) (http.Handler, error) {
	v := &v1{
		db: db,
	}
	v.initSchemaTypes()
	schema, err := v.getSchema()
	if err != nil {
		return nil, err
	}
	h := handler.New(&handler.Config{
		Schema: schema,
		Pretty: true,
	})
	return &graphqlHandlerWrapper{h}, nil
}
