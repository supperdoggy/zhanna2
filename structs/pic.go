package structs

import "time"

type Pic struct {
	ID string `json:"id" bson:"_id"`
	Data string `json:"data" bson:"data"`
	TimeCreated time.Time `json:"time_created" bson:"time_created"`
}
