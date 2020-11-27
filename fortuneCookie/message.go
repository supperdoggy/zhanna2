package main

import (
	"strings"
)

var (
	sign = []string{
		"овен",
		"телец",
		"близнецы",
		"рак",
		"лев",
		"дева",
		"весы",
		"скорпион",
		"стрелец",
		"козерог",
		"водолей",
		"рыбы",
	}
)

type DownloadedData struct {
	Name string `json:"name"`
	Messages []Message `json:"messages"`
}

type Cookie struct {
	Id int `json:"_id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}

type Message struct {
	ID   uint64   `json:"id"`
	Type string   `json:"type"`
	Text string `json:"text"`
}

func (m *Message) checkIfFortune() bool {
	if len(m.Text) < 8 {
		return false
	}
	if stringInArr(strings.ToLower(m.Text[0:CharLocation(m.Text, ".", 1)]), sign) {
		return true
	}
	return false
}

func CharLocation(s, c string, amount int) int{
	am := 0
	for k, v := range s{
		if string(v) == c{
			am++
			if am == amount{
				return k
			}
		}
	}
	return 0
}

func stringInArr(s string, ar []string) bool {
	for _, v := range ar {
		if v == s {
			return true
		}
	}
	return false
}
