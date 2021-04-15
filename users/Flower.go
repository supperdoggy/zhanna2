package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Flower - struct of flower
type Flower struct {
	ID   uint64 `json:"id" bson:"_id"`
	HP   uint8  `json:"hp" bson:"hp"`
	Name string `json:"name" bson:"name"`
	Icon string `json:"icon" bson:"icon"`
	Type string `json:"type" bson:"type"`
	// owner id
	Owner int   `json:"owner" bson:"owner"`
	Grew  uint8 `json:"grew" bson:"grew"`
	Dead  bool  `json:"dead" bson:"dead"`

	CreationTime time.Time `json:"creationTime" bson:"creationTime"`
	LastTimeGrow time.Time `json:"lastTimeGrow" bson:"lastTimeGrow"`
}

func canGrowFlower(id int) (bool, error) {
	answer, err := MakeReqToFlowers("canGrowFlower", obj{"id": id})
	if err != nil {
		fmt.Println("canGrowFlower() -> MakeReqToFlower(canGrowFlower) error:", err.Error())
		return false, err
	}

	var answerStruct struct {
		Answer bool   `json:"answer"`
		Err    string `json:"err"`
	}
	if err := json.Unmarshal(answer, &answerStruct); err != nil {
		fmt.Println("canGrowFlower() -> Unmarshal error:", err.Error(), string(answer))
		return false, err
	}

	if answerStruct.Err != "" {
		fmt.Println("canGrowFlower() -> got error from flower:", answerStruct.Err)
		return false, fmt.Errorf(answerStruct.Err)
	}
	return answerStruct.Answer, nil

}
