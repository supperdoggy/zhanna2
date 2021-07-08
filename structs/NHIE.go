package structs

type NHIE struct {
	ID   int    `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}