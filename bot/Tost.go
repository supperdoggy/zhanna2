package main

type Tost struct {
	ID   int    `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}