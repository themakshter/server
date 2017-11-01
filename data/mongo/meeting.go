package mongo

import (
	"errors"
	"time"

	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/auth"
	"github.com/impactasaurus/server/data"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func enforceMeetingPermissions(u auth.User) (query, error) {
	if u.IsBeneficiary() {
		meeting, ok := u.GetAssessmentScope()
		if !ok {
			return nil, errors.New("Not authorized to access assessment")
		}
		return query{
			"_id":         meeting,
			"beneficiary": u.UserID(),
		}, nil
	}
	userOrg, err := u.Organisation()
	if err != nil {
		return nil, err
	}
	return query{
		"organisationID": userOrg,
	}, nil
}

func (m *mongo) GetMeeting(id string, u auth.User) (impact.Meeting, error) {
	meeting := impact.Meeting{}

	col, closer := m.getMeetingCollection()
	defer closer()

	query, err := enforceMeetingPermissions(u)
	if err != nil {
		return impact.Meeting{}, err
	}
	if err = query.addQueryFields(map[string]interface{}{
		"_id": id,
	}); err != nil {
		return impact.Meeting{}, err
	}

	err = col.Find(query).One(&meeting)
	if err != nil {
		if mgo.ErrNotFound == err {
			return meeting, data.NewNotFoundError("Meeting")
		}
		return meeting, err
	}
	return meeting, nil
}

type meetingGetter func(col *mgo.Collection, q query) ([]impact.Meeting, error)

func (m *mongo) getMeetings(inner meetingGetter, u auth.User) ([]impact.Meeting, error) {
	col, closer := m.getMeetingCollection()
	defer closer()

	baseQuery, err := enforceMeetingPermissions(u)
	if err != nil {
		return nil, err
	}
	return inner(col, baseQuery)
}

func (m *mongo) GetMeetingsForBeneficiary(beneficiary string, u auth.User) ([]impact.Meeting, error) {
	return m.getMeetings(func(col *mgo.Collection, q query) ([]impact.Meeting, error) {
		if err := q.addQueryFields(map[string]interface{}{
			"beneficiary": beneficiary,
		}); err != nil {
			return nil, err
		}
		results := []impact.Meeting{}
		err := col.Find(q).All(&results)
		return results, err
	}, u)
}

func (m *mongo) GetOSMeetingsForBeneficiary(beneficiary string, outcomeSetID string, u auth.User) ([]impact.Meeting, error) {
	return m.getMeetings(func(col *mgo.Collection, q query) ([]impact.Meeting, error) {
		if err := q.addQueryFields(map[string]interface{}{
			"beneficiary":  beneficiary,
			"outcomeSetID": outcomeSetID,
		}); err != nil {
			return nil, err
		}
		results := []impact.Meeting{}
		err := col.Find(q).All(&results)
		return results, err
	}, u)
}

func (m *mongo) GetOSMeetingsInTimeRange(start, end time.Time, outcomeSetID string, u auth.User) ([]impact.Meeting, error) {
	return m.getMeetings(func(col *mgo.Collection, q query) ([]impact.Meeting, error) {
		if err := q.addQueryFields(map[string]interface{}{
			"outcomeSetID": outcomeSetID,
			"conducted": bson.M{
				"$gte": start,
				"$lte": end,
			},
		}); err != nil {
			return nil, err
		}
		results := []impact.Meeting{}
		err := col.Find(q).All(&results)
		return results, err
	}, u)
}

func (m *mongo) NewMeeting(beneficiaryID, outcomeSetID string, conducted time.Time, u auth.User) (impact.Meeting, error) {
	userOrg, err := u.Organisation()
	if err != nil {
		return impact.Meeting{}, err
	}

	col, closer := m.getMeetingCollection()
	defer closer()

	meeting := impact.Meeting{
		ID:             uuid.NewV4().String(),
		OrganisationID: userOrg,
		OutcomeSetID:   outcomeSetID,
		Beneficiary:    beneficiaryID,
		Conducted:      conducted,
		Created:        time.Now(),
		Modified:       time.Now(),
		User:           u.UserID(),
	}

	if err := col.Insert(meeting); err != nil {
		return impact.Meeting{}, err
	}
	return meeting, nil
}

func (m *mongo) NewAnswer(meetingID string, answer impact.Answer, u auth.User) (impact.Meeting, error) {
	q, err := enforceMeetingPermissions(u)
	if err != nil {
		return impact.Meeting{}, err
	}

	col, closer := m.getMeetingCollection()
	defer closer()

	if err := q.addQueryFields(map[string]interface{}{
		"_id": meetingID,
	}); err != nil {
		return impact.Meeting{}, err
	}

	if err := col.Update(q, bson.M{
		"$push": bson.M{
			"answers": answer,
		},
	}); err != nil {
		return impact.Meeting{}, err
	}

	return m.GetMeeting(meetingID, u)
}
