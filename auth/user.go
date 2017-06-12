package auth

type User interface {
	Organisation() (string, error)
}

type auth0User struct {
}

func newUser(jwt string) (User, error) {
	return &auth0User{}, nil
}

func (u *auth0User) Organisation() (string, error) {
	return "test org", nil
}
