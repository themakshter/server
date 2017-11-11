package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	jwtLib "github.com/dgrijalva/jwt-go"
	"github.com/impactasaurus/server/auth"
	"github.com/stretchr/testify/assert"
)

const aud = "test-aud"
const iss = "test-iss"

func getTarget(t *testing.T) (*rsa.PublicKey, auth.Generator, auth.Authenticator) {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	assert.Nil(t, err)
	return &key.PublicKey, auth.NewBeneficiaryJWTGenerator(aud, iss, key), auth.NewJWTAuthenticator(aud, iss, &key.PublicKey)
}

func TestBeneficiaryJWTContent(t *testing.T) {
	pubKey, target, authenticator := getTarget(t)
	meeting := "m1"
	benID := "ben1"
	jti, token, err := target.GenerateBeneficiaryJWT(benID, meeting, time.Minute)
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

	parsedToken, err := jwtLib.ParseWithClaims(token, &jwtLib.StandardClaims{}, func(token *jwtLib.Token) (interface{}, error) {
		return pubKey, nil
	})
	assert.Nil(t, err)

	claims, ok := parsedToken.Claims.(*jwtLib.StandardClaims)
	assert.True(t, ok)
	assert.Equal(t, jti, claims.Id)
}

func TestExpiry(t *testing.T) {
	_, target, authenticator := getTarget(t)
	meeting := "m1"
	benID := "ben1"
	_, token, err := target.GenerateBeneficiaryJWT(benID, meeting, time.Second)
	assert.Nil(t, err)
	time.Sleep(time.Second * 2)
	_, err = authenticator.AuthUser(token)
	assert.NotNil(t, err)
}
