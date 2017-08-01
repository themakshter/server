package server

const LIKERT = "likert"

const (
	MEAN = "mean"
	SUM  = "sum"
)

type Question struct {
	ID         string                 `json:"id"`
	Question   string                 `json:"question"`
	Type       string                 `json:"type"`
	Deleted    bool                   `json:"deleted"`
	Options    map[string]interface{} `json:"options"`
	CategoryID string                 `json:"categoryID"  bson:"categoryID"`
}

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Aggregation string `json:"aggregation"`
}

type OutcomeSet struct {
	ID             string     `json:"id" bson:"_id"`
	OrganisationID string     `json:"organisationID" bson:"organisationID"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Questions      []Question `json:"questions"`
	Categories     []Category `json:"categories"`
	Deleted        bool       `json:"deleted"`
}

func (os *OutcomeSet) GetCategoryQuestions(catID string) []Question {
	out := make([]Question, 0, len(os.Questions))
	for _, q := range os.Questions {
		if q.CategoryID == catID {
			out = append(out, q)
		}
	}
	return out
}
