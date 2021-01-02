package main

type Tost struct {
	ID   int    `bson:"_id" json:"id"`
	Text string `bspn:"text" json:"text"`
}
