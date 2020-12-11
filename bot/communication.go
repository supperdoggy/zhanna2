package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// MakeRandomAnekHttpReq - sends http req to anek server and unmarshals it to RandomAnekAnswer struct
func MakeRandomAnekHttpReq()(response RandomAnekAnswer, err error){
	resp, err := MakeHttpReq(anekUrl+"/api/v1/getRandomAnek", "GET", nil)
	if err != nil {
		fmt.Println("handlers.go -> MakeRandomAnekHttpReq() -> MakeHttpReq ->", err.Error())
		return
	}

	if err = json.Unmarshal(resp, &response);err!=nil{
		fmt.Println("communication.go -> MakeRandomAnekHttpReq() -> error ->", err.Error())
		return
	}

	return
}

// MakeHttpReq - func for sending http req with given path, method(get or post!) and data
func MakeHttpReq(path, method string, data io.Reader) (answer []byte, err error) {
	var resp *http.Response
	switch method {
	case "GET":
		resp, err = http.Get(path)
	case "POST":
		resp, err = http.Post(path, "application/json", data)
	default:
		err = fmt.Errorf("method not supported, use get or post methods")
	}
	if err != nil {
		return nil, err
	}

	answer, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}
