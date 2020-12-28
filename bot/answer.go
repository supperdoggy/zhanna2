package main

type RandomAnekAnswer struct {
	Anek
	Err string `json:"err"`
}

type Anek struct {
	Id   int    `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}
