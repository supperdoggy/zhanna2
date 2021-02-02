package main

import "time"

// Flower - struct of flower
type Flower struct {
	ID   uint64 `json:"id" bson:"_id"`
	HP   uint8  `json:"hp" bson:"hp,omitempty"`
	Name string `json:"name" bson:"name"`
	Icon string `json:"icon" bson:"icon"`
	Type string `json:"type" bson:"type"`
	// owner id
	Owner int   `json:"owner" bson:"owner"`
	Grew  uint8 `json:"grew" bson:"grew"`

	CreationTime time.Time `json:"creationTime" bson:"creationTime"`
	LastTimeGrow time.Time `json:"lastTimeGrow" bson:"lastTimeGrow"`
}