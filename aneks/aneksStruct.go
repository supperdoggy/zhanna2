package main

type Aneks struct {
	A []Anek `json:"anek" bson:"anek"`
}

type Anek struct {
	Id   int    `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}


