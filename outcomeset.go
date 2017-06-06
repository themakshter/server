package server

const LIKERT = "likert"

type Question struct {
	ID       string                 `json:"id"`
	Question string                 `json:"question"`
	Type     string                 `json:"type"`
	Deleted  bool                   `json:"deleted"`
	Options  map[string]interface{} `json:"options"`
}

type OutcomeSet struct {
	ID             string     `json:"id" bson:"_id"`
	OrganisationID string     `json:"organisationID" bson:"organisationID"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Questions      []Question `json:"questions"`
	Deleted        bool       `json:"deleted"`
}
