package api

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/impactasaurus/server/data"
	"net/http"
)

type Listener interface {
	Listen() error
}

type v1 struct {
	db                  data.Base
	questionType        *graphql.Object
	questionInputType   *graphql.InputObject
	outcomeSetType      *graphql.Object
	outcomeSetInputType *graphql.InputObject
	answerType          *graphql.Object
	organisationType    *graphql.Object
	meetingType         *graphql.Object
}

func NewV1(db data.Base) Listener {
	v := &v1{
		db: db,
	}
	v.initSchemaTypes()
	return v
}

func (v *v1) Listen() error {
	schema, err := v.getSchema()
	if err != nil {
		return err
	}
	h := handler.New(&handler.Config{
		Schema: schema,
		Pretty: true,
	})
	http.Handle("/graphql", h)
	return nil
}
