package main

import (
	"math/rand"
	"time"
)

func calculateChance(chance int) bool {
	rand.Seed(time.Now().Unix())
	var num int = rand.Intn(101)
	return num <= chance
}
