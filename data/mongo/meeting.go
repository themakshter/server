package mongo

import (
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (m *mongo) GetMeeting(id string, u auth.User) (*impact.Meeting, error) {
	col, closer := m.getMeetingCollection()
	defer closer()

	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}

	meeting := &impact.Meeting{}
	err = col.Find(bson.M{
		"_id": id,
		"organisationID": userOrg,
	}).One(meeting)
	if err != nil {
		if mgo.ErrNotFound == err {
			return nil, &data.NotFound{}
		}
		return nil, err
	}
	return meeting, nil
}

func (m *mongo) GetMeetingsForBeneficiary(beneficiary string, u auth.User) ([]impact.Meeting, error) {
	col, closer := m.getMeetingCollection()
	defer closer()

	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}

	results := []impact.Meeting{}
	err = col.Find(bson.M{
		"beneficiary": beneficiary,
		"organisationID": userOrg,
	}).All(&results)
	return results, err
}
