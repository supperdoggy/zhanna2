package structs

type Anek struct {
	Id   int    `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}
