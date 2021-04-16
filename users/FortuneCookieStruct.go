package main

import "time"

type FortuneCookie struct {
	Id   int    `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}

func CanGetFortune(date time.Time) bool {
	now := time.Now()
	return date.Day != now.Day || date.Month != now.Month || date.Year != now.Year
}
