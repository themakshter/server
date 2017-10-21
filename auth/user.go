package auth

// User is an object which provides details about the user making the request to the API
type User interface {
	// Organisation gets the organisation the user belongs to
	// errors are expected if the user is a beneficiary
	Organisation() (string, error)
	// UserID gets the user's ID within the system
	// for users this will be their auth0 IDs
	// for beneficiaries this will be their beneficiary ID
	UserID() string
	// IsBeneficiary returns true if the User is a beneficiary user
	// beneficiary users do not belong to an organisation and are normally limited in scope
	IsBeneficiary() bool
	// GetAssessmentScope returns true and the assessment ID if the user is restricted in scope to a single assessment
	// this is common for beneficiary users
	GetAssessmentScope() (string, bool)
}

// Authenticator takes a JWT, validates the JWT and generates a User object
type Authenticator interface {
	AuthUser(jwt string) (User, error)
}
