package auth

// User is an object which provides details about the user making the request to the API
type User interface {
	Organisation() (string, error)
	UserID() string
}

// Authenticator takes a JWT, validates the JWT and generates a User object
type Authenticator interface {
	AuthUser(jwt string) (User, error)
}
