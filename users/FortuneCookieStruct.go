package main

import "time"

type FortuneCookie struct {
	Id   int    `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}

func CanGetFortune(date time.Time) bool {
	return int(time.Now().Sub(date).Hours())/24 >= 1
}
