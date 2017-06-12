package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		authString := req.Header.Get("Authorization")
		fmt.Println(authString)
		if strings.HasPrefix(authString, "Bearer ") {
			jwt := strings.TrimPrefix(authString, "Bearer ")
			user, err := newUser(jwt)
			if err != nil {
				fmt.Println("error: user create failure")
				ctx = newContextWithAuthError(ctx, err)
			} else {
				ctx = newContextWithUser(ctx, user)
				fmt.Println("success")
				u, e := GetUser(ctx)
				fmt.Println(u)
				fmt.Println(e)
				fmt.Println(u.Organisation())
			}
		} else if authString == "" {
			fmt.Println("error: no auth header")
			ctx = newContextWithAuthError(ctx, errors.New("Missing Authorization header"))
		} else {
			fmt.Println("error: invalid header")
			ctx = newContextWithAuthError(ctx, errors.New("Invalid Authorization header"))
		}

		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}
