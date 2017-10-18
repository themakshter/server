package auth

import (
	"errors"
	"net/http"
	"strings"
)

func getUser(jwt string, auth ...Authenticator) (User, error) {
	for _, a := range auth {
		if u, err := a.AuthUser(jwt); err == nil {
			return u, err
		}
	}
	return nil, errors.New("Failed to authenticate user")
}

// Middleware creates a http handler middleware which authenticates responses using the provided Authenticators
// If authentication suceeds, the request will have a context which includes a User object
// If authentication fails, the context will have an authentication error
func Middleware(next http.Handler, auth ...Authenticator) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		authString := req.Header.Get("Authorization")
		if strings.HasPrefix(authString, "Bearer ") {
			jwt := strings.TrimPrefix(authString, "Bearer ")
			user, err := getUser(jwt, auth...)
			if err != nil {
				ctx = newContextWithAuthError(ctx, err)
			} else {
				ctx = newContextWithUser(ctx, user)
			}
		} else if authString == "" {
			ctx = newContextWithAuthError(ctx, errors.New("Missing Authorization header"))
		} else {
			ctx = newContextWithAuthError(ctx, errors.New("Invalid Authorization header"))
		}

		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}
