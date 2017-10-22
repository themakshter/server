package auth

import (
	"crypto/rsa"

	jwtLib "github.com/dgrijalva/jwt-go"
	"github.com/impactasaurus/server/log"
)

func MustParseRSAPublicKeyFromPEM(key string) *rsa.PublicKey {
	pubKey, err := jwtLib.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		log.Fatal(err, nil)
	}
	return pubKey
}

func MustParseRSAPrivateKeyFromPEM(key string) *rsa.PrivateKey {
	privateKey, err := jwtLib.ParseRSAPrivateKeyFromPEM([]byte(key))
	if err != nil {
		log.Fatal(err, nil)
	}
	return privateKey
}
