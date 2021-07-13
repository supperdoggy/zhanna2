package structs

import (
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

type Message struct {
	UserID  int             `json:"userID" bson:"userID"`
	Message telebot.Message `json:"message" bson:"message"`
	Reply   telebot.Message `json:"reply" bson:"reply"`
	Time    time.Time       `json:"time" bson:"time"`
}
