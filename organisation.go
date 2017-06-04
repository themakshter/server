package server

type Organisation struct {
	Name string `json:"name"`
	ID   string `json:"id",bson:"_id"`
}
