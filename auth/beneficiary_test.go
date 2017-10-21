package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/impactasaurus/server/auth"
	"github.com/stretchr/testify/assert"
)

const aud = "test-aud"
const iss = "test-iss"

func getTarget(t *testing.T) (auth.Generator, auth.Authenticator) {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	assert.Nil(t, err)
	return auth.NewBeneficiaryJWTGenerator(aud, iss, key), auth.NewJWTAuthenticator(aud, iss, &key.PublicKey)
}

func TestBeneficiaryJWTContent(t *testing.T) {
	target, authenticator := getTarget(t)
	meeting := "m1"
	benID := "ben1"
	token, err := target.GenerateBeneficiaryJWT(benID, meeting, time.Minute)
	assert.Nil(t, err)
	u, err := authenticator.AuthUser(token)
	assert.Nil(t, err)
	assert.Equal(t, u.UserID(), benID)
	assert.Equal(t, u.IsBeneficiary(), true)
	m, ok := u.GetAssessmentScope()
	assert.Equal(t, ok, true)
	assert.Equal(t, m, meeting)
	_, err = u.Organisation()
	assert.NotNil(t, err)
}

func TestExpiry(t *testing.T) {
	target, authenticator := getTarget(t)
	meeting := "m1"
	benID := "ben1"
	token, err := target.GenerateBeneficiaryJWT(benID, meeting, time.Second)
	assert.Nil(t, err)
	time.Sleep(time.Second * 2)
	_, err = authenticator.AuthUser(token)
	assert.NotNil(t, err)
}
