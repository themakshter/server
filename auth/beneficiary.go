package auth

import (
	"crypto/rsa"
	"time"

	"errors"

	jwtLib "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
)

const beneficiaryKey = "beneficiary"
const assessmentScopeKey = "scope"

type benGen struct {
	private *rsa.PrivateKey
	aud     string
	iss     string
}

// NewBeneficiaryJWTGenerator returns a beneficiary JWT generator using the provided
// audience, issuer and private RSA key
func NewBeneficiaryJWTGenerator(aud, iss string, private *rsa.PrivateKey) Generator {
	return &benGen{
		private: private,
		aud:     aud,
		iss:     iss,
	}
}

func (b *benGen) GenerateBeneficiaryJWT(benID, meetingID string, expiry time.Duration) (jti, jwt string, err error) {
	meta := map[string]interface{}{}
	meta[beneficiaryKey] = true
	meta[assessmentScopeKey] = meetingID
	jti = uuid.NewV4().String()
	token := jwtLib.NewWithClaims(jwtLib.SigningMethodRS256, jwtUser{
		AppMetadata: meta,
		StandardClaims: jwtLib.StandardClaims{
			Audience:  b.aud,
			ExpiresAt: time.Now().Add(expiry).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    b.iss,
			Subject:   benID,
			Id:        jti,
		},
	})
	jwt, err = token.SignedString(b.private)
	return
}

type benAuth struct {
	inner Authenticator
}

// NewBeneficiaryAuthenticator returns an Authenticator which authenticates only beneficiary JWTs
func NewBeneficiaryAuthenticator(aud, iss string, key *rsa.PublicKey) Authenticator {
	return &benAuth{
		inner: NewJWTAuthenticator(aud, iss, key),
	}
}

func (b *benAuth) AuthUser(jwt string) (User, error) {
	u, err := b.inner.AuthUser(jwt)
	if err != nil {
		return u, err
	}
	if ben := u.IsBeneficiary(); ben == false {
		return nil, errors.New("JWT was not a beneficiary JWT")
	}
	return u, err
}
