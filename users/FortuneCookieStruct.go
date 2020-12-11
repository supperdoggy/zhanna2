package main

type FortuneCookie struct {
	Id int `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}
