// +build integration

package mongo_test

import (
	"testing"
	"time"

	"errors"

	"github.com/golang/mock/gomock"
	impact "github.com/impactasaurus/server"
	"github.com/impactasaurus/server/data"
	"github.com/impactasaurus/server/data/mongo"
	"github.com/impactasaurus/server/mock"
	"github.com/kelseyhightower/envconfig"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func seedMeeting(t *testing.T, mockCtrl *gomock.Controller, m data.Base, orgID, userID, benID, osID string, conducted time.Time) impact.Meeting {
	mockUser := mock.NewMockUser(mockCtrl)
	mockUser.EXPECT().Organisation().Return(orgID, nil)
	mockUser.EXPECT().UserID().Return(userID)
	meeting, err := m.NewMeeting(benID, osID, conducted, mockUser)
	require.NoError(t, err)
	return meeting
}

func getTarget(t *testing.T) data.Base {
	c := &mongo.Config{}
	envconfig.MustProcess("MONGO", c)
	m, err := mongo.New(c.URL, c.Port, uuid.NewV4().String(), c.User, c.Password)
	require.NoError(t, err)
	return m
}

func TestGetMeeting(t *testing.T) {
	target := getTarget(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	os1Meeting := seedMeeting(t, mockCtrl, target, "org1", "user1", "ben1", "os1", time.Now())
	os2Meeting := seedMeeting(t, mockCtrl, target, "org1", "user2", "ben1", "os2", time.Now())
	seedMeeting(t, mockCtrl, target, "org1", "user2", "ben2", "os3", time.Now())

	// organisation user searching for meeting
	mockUser := mock.NewMockUser(mockCtrl)
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org1", nil)
	meeting, err := target.GetMeeting(os1Meeting.ID, mockUser)
	require.NoError(t, err)
	require.Equal(t, "user1", meeting.User)

	// Non beneficiary user without an organisation
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org1", errors.New("err"))
	meeting, err = target.GetMeeting(os1Meeting.ID, mockUser)
	require.Error(t, err)

	// beneficiary searching for their scoped meeting
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben1")
	meeting, err = target.GetMeeting(os1Meeting.ID, mockUser)
	require.NoError(t, err)
	require.Equal(t, "user1", meeting.User)

	// beneficiary user searching for another beneficiary's meeting
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben2")
	meeting, err = target.GetMeeting(os1Meeting.ID, mockUser)
	require.Error(t, err)

	// beneficiary user searching for their non scoped meeting
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben1")
	meeting, err = target.GetMeeting(os2Meeting.ID, mockUser)
	require.Error(t, err)

	// beneficiary without an assessment scope
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return("", false)
	meeting, err = target.GetMeeting(os1Meeting.ID, mockUser)
	require.Error(t, err)

	// beneficiary with unknown assessment scope
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return("unknown", true)
	mockUser.EXPECT().UserID().Return("ben1")
	meeting, err = target.GetMeeting(os1Meeting.ID, mockUser)
	require.Error(t, err)
}

func TestGetMeetingForBen(t *testing.T) {
	target := getTarget(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	os1Meeting := seedMeeting(t, mockCtrl, target, "org1", "user1", "ben1", "os1", time.Now())
	seedMeeting(t, mockCtrl, target, "org1", "user2", "ben1", "os2", time.Now())
	// target beneficiary in a different organisation
	seedMeeting(t, mockCtrl, target, "org2", "user3", "ben1", "os3", time.Now())
	// different beneficiary in user's organisation
	seedMeeting(t, mockCtrl, target, "org1", "user4", "ben2", "os4", time.Now())

	// organisation user searching for beneficiary
	mockUser := mock.NewMockUser(mockCtrl)
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org1", nil)
	meetings, err := target.GetMeetingsForBeneficiary("ben1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 2)
	require.Equal(t, "os1", meetings[0].OutcomeSetID)
	require.Equal(t, "os2", meetings[1].OutcomeSetID)

	// Non beneficiary user without an organisation
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org2", errors.New("err"))
	meetings, err = target.GetMeetingsForBeneficiary("ben1", mockUser)
	require.Error(t, err)

	// beneficiary searching for their scoped meeting
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben1")
	meetings, err = target.GetMeetingsForBeneficiary("ben1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 1)
	require.Equal(t, "os1", meetings[0].OutcomeSetID)

	// beneficiary user searching for another beneficiary
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben2")
	meetings, err = target.GetMeetingsForBeneficiary("ben1", mockUser)
	require.Error(t, err)

	// beneficiary without an assessment scope
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return("", false)
	meetings, err = target.GetMeetingsForBeneficiary("ben1", mockUser)
	require.Error(t, err)

	// beneficiary with unknown assessment scope
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return("unknown", true)
	mockUser.EXPECT().UserID().Return("ben1")
	meetings, err = target.GetMeetingsForBeneficiary("ben1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 0)
}

func TestGetOSMeetingForBen(t *testing.T) {
	target := getTarget(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	os1Meeting := seedMeeting(t, mockCtrl, target, "org1", "user1", "ben1", "os1", time.Now())
	seedMeeting(t, mockCtrl, target, "org1", "user2", "ben1", "os2", time.Now())
	// target beneficiary in a different organisation
	seedMeeting(t, mockCtrl, target, "org2", "user3", "ben1", "os3", time.Now())
	// different beneficiary in user's organisation
	seedMeeting(t, mockCtrl, target, "org1", "user4", "ben2", "os4", time.Now())

	// organisation user searching for beneficiary and OS
	mockUser := mock.NewMockUser(mockCtrl)
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org1", nil)
	meetings, err := target.GetOSMeetingsForBeneficiary("ben1", "os1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 1)
	require.Equal(t, "os1", meetings[0].OutcomeSetID)

	// Non beneficiary user without an organisation
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org2", errors.New("err"))
	meetings, err = target.GetOSMeetingsForBeneficiary("ben1", "os1", mockUser)
	require.Error(t, err)

	// beneficiary searching for their scoped meeting
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben1")
	meetings, err = target.GetOSMeetingsForBeneficiary("ben1", "os1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 1)
	require.Equal(t, "os1", meetings[0].OutcomeSetID)

	// beneficiary user searching for another beneficiary
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben2")
	meetings, err = target.GetOSMeetingsForBeneficiary("ben1", "os1", mockUser)
	require.Error(t, err)

	// beneficiary without an assessment scope
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return("", false)
	meetings, err = target.GetOSMeetingsForBeneficiary("ben1", "os1", mockUser)
	require.Error(t, err)

	// beneficiary with unknown assessment scope
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return("unknown", true)
	mockUser.EXPECT().UserID().Return("ben1")
	meetings, err = target.GetOSMeetingsForBeneficiary("ben1", "os1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 0)
}

func TestGetOSMeetingInTimeRange(t *testing.T) {
	target := getTarget(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	os1Meeting := seedMeeting(t, mockCtrl, target, "org1", "user1", "ben1", "os1", time.Unix(1, 0))
	seedMeeting(t, mockCtrl, target, "org1", "user2", "ben2", "os1", time.Unix(2, 0))
	seedMeeting(t, mockCtrl, target, "org2", "user3", "ben3", "os2", time.Unix(3, 0))

	// organisation user searching for OS over a range containing all our seeded meetings
	mockUser := mock.NewMockUser(mockCtrl)
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org1", nil)
	meetings, err := target.GetOSMeetingsInTimeRange(time.Unix(0, 0), time.Unix(4, 0), "os1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 2)
	require.Equal(t, "user1", meetings[0].User)
	require.Equal(t, "user2", meetings[1].User)

	// Non beneficiary user without an organisation
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org2", errors.New("err"))
	meetings, err = target.GetOSMeetingsInTimeRange(time.Unix(0, 0), time.Unix(4, 0), "os1", mockUser)
	require.Error(t, err)

	// organisation user searching for OS over sub range
	mockUser.EXPECT().IsBeneficiary().Return(false)
	mockUser.EXPECT().Organisation().Return("org1", nil)
	meetings, err = target.GetOSMeetingsInTimeRange(time.Unix(0, 0), time.Unix(1, 0), "os1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 1)
	require.Equal(t, "user1", meetings[0].User)

	// beneficiary searching over all the range including their scoped meeting
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben1")
	meetings, err = target.GetOSMeetingsInTimeRange(time.Unix(0, 0), time.Unix(5, 0), "os1", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 1)
	require.Equal(t, "os1", meetings[0].OutcomeSetID)

	// beneficiary searching for outcome set which their scoped meeting is not part of
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return(os1Meeting.ID, true)
	mockUser.EXPECT().UserID().Return("ben1")
	meetings, err = target.GetOSMeetingsInTimeRange(time.Unix(0, 0), time.Unix(5, 0), "os2", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 0)

	// beneficiary without an assessment scope
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return("", false)
	meetings, err = target.GetOSMeetingsInTimeRange(time.Unix(0, 0), time.Unix(5, 0), "os2", mockUser)
	require.Error(t, err)

	// beneficiary with unknown assessment scope
	mockUser.EXPECT().IsBeneficiary().Return(true)
	mockUser.EXPECT().GetAssessmentScope().Return("unknown", true)
	mockUser.EXPECT().UserID().Return("ben1")
	meetings, err = target.GetOSMeetingsInTimeRange(time.Unix(0, 0), time.Unix(5, 0), "os2", mockUser)
	require.NoError(t, err)
	require.Len(t, meetings, 0)
}
