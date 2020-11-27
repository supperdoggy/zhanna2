package main
//
//import (
//	"encoding/json"
//	"fmt"
//	"gopkg.in/mgo.v2"
//	"io/ioutil"
//)
//
//type obj map[string]interface{}
//type arr []interface{}
//
//func connectToCookiesDb() *mgo.Collection {
//	m, err := mgo.Dial("")
//	if err != nil {
//		panic(err.Error())
//	}
//	cookiesCollection := m.DB("fortuneCookie").C("Cookies")
//	return cookiesCollection
//}
//
//func getDownloadedJsonData() (result DownloadedData) {
//	d, err := ioutil.ReadFile("/Users/maks/go/src/github.com/supperdoggy/superSecretDevelopement/fortuneCookie/result.json")
//	if err != nil {
//		panic(err.Error())
//	}
//
//	if err = json.Unmarshal(d, &result); err != nil {
//
//	}
//
//	return
//}
//
////func convertMapToStruct(data []map[string]interface{}) (result []Message) {
////	for k, v := range data {
////		result = append(result, Message{
////			ID:   v["id"].(uint64),
////			Type: v["type"].(string),
////			Text: v["text"].([]string),
////		})
////		fmt.Println(k)
////	}
////	return
////}
//
//func main() {
//	jsonCookies := getDownloadedJsonData()
//	count := 0
//	messageSlice := []Message{}
//	for k, _ := range jsonCookies.Messages{
//		if jsonCookies.Messages[k].checkIfFortune(){
//			count++
//			messageSlice = append(messageSlice, jsonCookies.Messages[k])
//		}
//	}
//	fmt.Println(count)
//	col := connectToCookiesDb()
//	cookies := []Cookie{}
//	for k, _ := range messageSlice{
//		cookies = append(cookies, Cookie{
//			Id:   k,
//			Text: messageSlice[k].Text[CharLocation(messageSlice[k].Text, ".", 2)+1:],
//		})
//		if err := col.Insert(Cookie{
//			Id:   k,
//			Text: messageSlice[k].Text[CharLocation(messageSlice[k].Text, ".", 2)+1:],
//		}); err != nil{
//			fmt.Println(k, "error")
//		}
//
//	}
//	fmt.Println("Done inserting!")
//
//}
