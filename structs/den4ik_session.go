package structs

import "time"

type Session struct {
	ID           int       `json:"id" bson:"_id"` // chat_id
	Cards        []Card  `json:"cards" bson:"cards"`
	CreationTime time.Time `json:"creation_time" bson:"creation_time"`
}
