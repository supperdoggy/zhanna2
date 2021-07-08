package structs

import (
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"math/rand"
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

func FlowerDies() bool {
	rand.Seed(time.Now().Unix())
	var num int = rand.Intn(101)
	return num <= cfg.DyingChance
}
