package main

type Cookie struct {
	ID   int32  `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}
