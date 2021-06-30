package structs

type Tost struct {
	ID   int    `bson:"_id" json:"id"`
	Text string `bson:"text" json:"text"`
}
