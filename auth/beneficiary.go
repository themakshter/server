package auth

import (
	"crypto/rsa"
	"time"

	jwtLib "github.com/dgrijalva/jwt-go"
)

const beneficiaryKey = "beneficiary"
const assessmentScopeKey = "scope"

// Generator generates beneficiary JWTs for use with self assessments
type Generator interface {
	GenerateBeneficiaryJWT(benID, meetingID string, expiry time.Duration) (string, error)
}

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

func (b *benGen) GenerateBeneficiaryJWT(benID, meetingID string, expiry time.Duration) (string, error) {
	meta := map[string]interface{}{}
	meta[beneficiaryKey] = true
	meta[assessmentScopeKey] = meetingID
	token := jwtLib.NewWithClaims(jwtLib.SigningMethodRS256, jwtUser{
		AppMetadata: meta,
		StandardClaims: jwtLib.StandardClaims{
			Audience:  b.aud,
			ExpiresAt: time.Now().Add(expiry).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    b.iss,
			Subject:   benID,
		},
	})
	return token.SignedString(b.private)
}
