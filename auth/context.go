package auth

import (
	"context"
	"errors"
	"fmt"
)

type key int

const userIndex key = 0
const authErrorIndex key = 1

func newContextWithUser(ctx context.Context, userObj User) context.Context {
	return context.WithValue(ctx, userIndex, userObj)
}

func newContextWithAuthError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, authErrorIndex, err)
}

func GetUser(ctx context.Context) (User, error) {
	user, uOk := ctx.Value(userIndex).(User)
	err, eOk := ctx.Value(authErrorIndex).(error)

	fmt.Println("getuser")
	fmt.Println(uOk)
	fmt.Println(eOk)
	fmt.Println(user)
	fmt.Println(err)

	if user != nil && uOk && err == nil {
		return user, nil
	}

	if err != nil && eOk {
		return nil, err
	}

	return nil, errors.New("Failed accessing user information from context")
}
