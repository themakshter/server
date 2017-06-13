package auth

import (
	"errors"
	"fmt"
	jwtLib "github.com/dgrijalva/jwt-go"
)

const aud = "pfKiAOUJh5r6jCxRn5vUYq7odQsjPUKf"
const iss = "https://impact.eu.auth0.com/"
const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA72wlFrwLpR2RV5+TxcQU
iLqAqILTRLLWkYJnw+xePiuSIutXARkOcyKGfsunZ/1xZLe5rRuP06M/RTKCNBrq
WwcSLEjL7Dh4JbrkDCGWx7YFcKzswZ3/3Gb2BtavpDhfg5RzaIBgnC+uyKAZvocA
+YoKT3SPHm79vnmvQvpheMqfuNJSmw0mCXMcwqjJUQ2GLQsWfI2qeVybxcsRDqp5
5kDGPVThrg8OGwJQDrHE4DopXWHWvUxKVS/e6eFN6qgibVP7vD3SnL0M7wgJDQuk
TzskiF5Zzsgc86b2P6kQ1H2ryKp1jNDjkBCpr3F7KfNek/ADrSZVpxFSg4cve7X3
+wIDAQAB
-----END PUBLIC KEY-----`

type User interface {
	Organisation() (string, error)
	UserID() string
}

type auth0User struct {
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
	jwtLib.StandardClaims
}

func validateAndParseJWT(jwt string) (User, error) {
	token, err := jwtLib.ParseWithClaims(jwt, &auth0User{}, func(token *jwtLib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtLib.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwtLib.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*auth0User)
	if !ok {
		return nil, errors.New("Failed to parse token to user object")
	}
	if !claims.VerifyAudience(aud, true) {
		return nil, errors.New("Token validation error: aud")
	}
	if !claims.VerifyIssuer(iss, true) {
		return nil, errors.New("Token validation error: iss")
	}
	return claims, nil
}

func newUser(jwt string) (User, error) {
	return validateAndParseJWT(jwt)
}

func (u *auth0User) Organisation() (string, error) {
	org, ok := u.AppMetadata["organisation"].(string)
	if !ok {
		return "", errors.New("Failed to extract organisation")
	}
	return org, nil
}

func (u *auth0User) UserID() string {
	return u.Subject
}
