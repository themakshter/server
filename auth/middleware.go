package auth

import (
	"errors"
	"net/http"
	"strings"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		authString := req.Header.Get("Authorization")
		if strings.HasPrefix(authString, "Bearer ") {
			jwt := strings.TrimPrefix(authString, "Bearer ")
			user, err := newUser(jwt)
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
