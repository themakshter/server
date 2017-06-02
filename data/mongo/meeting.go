package mongo

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) GetMeeting(id string) (*impact.Meeting, error) {
	col, closer := m.getMeetingCollection()
	defer closer()

	meeting := &impact.Meeting{}
	err := col.FindId(id).One(meeting)
	if err != nil {
		if mgo.ErrNotFound == err {
			return nil, &data.NotFound{}
		}
		return nil, err
	}
	return meeting, nil
}

func (m *mongo) GetMeetingsForBeneficiary(beneficiary string) ([]impact.Meeting, error) {
	col, closer := m.getMeetingCollection()
	defer closer()

	results := []impact.Meeting{}
	err := col.Find(bson.M{"beneficiary": beneficiary}).All(&results)
	return results, err
}
