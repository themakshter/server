package server

type Question struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Type     string `json:"type"`
}

type OutcomeSet struct {
	ID             string     `json:"id",bson:"_id"`
	OrganisationID string     `json:"organisationID"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Questions      []Question `json:"questions"`
	Deleted        bool       `json:"deleted"`
}
