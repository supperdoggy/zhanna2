package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// )

// func TransferToMgo() {
// 	data, err := ioutil.ReadFile("/Users/maks/go/src/github.com/supperdoggy/superSecretDevelopement/tost/Tosts.json")
// 	if err != nil {
// 		panic(err.Error)
// 	}
// 	marshaled := []string{}
// 	if err := json.Unmarshal(data, &marshaled); err != nil {
// 		panic(err.Error())
// 	}
// 	for k, v := range marshaled {
// 		a := anek{
// 			Id:   k,
// 			Text: v,
// 		}
// 		if err := aneks.TostCollection.Insert(a); err != nil {
// 			fmt.Println("error inserting", k, err.Error())
// 		}
// 	}
// 	fmt.Println("Done, inserted", len(marshaled), "docs")
// }
