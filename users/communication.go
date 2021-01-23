package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MakeReqToAnek - makes req to anek service
func MakeReqToAnek(method string, data []byte) (answer []byte, err error) {
	path := fmt.Sprintf("%s/%s", anekUrl, method)
	switch method {
	case "getRandomAnek":
		answer, err = MakeHttpReq(path, "GET", data)
	default:
		err = fmt.Errorf("no such method")
	}
	return
}

// MakeReqToFlowers - makes req to flowers service
// TODO: implement all the methods
func MakeReqToFlowers(method string, data []byte) (answer []byte, err error) {
	path := fmt.Sprintf("%s/%s", flowerUrl, method)
	switch method {
	case "addFlower":
		answer, err = MakeHttpReq(path, "POST", data)
	case "growFlower":
		answer, err = MakeHttpReq(path, "POST", data)
	case "canGrowFlower":
		answer, err = MakeHttpReq(path, "POST", data)
	case "getUserFlowers":
		answer, err = MakeHttpReq(path, "POST", data)
	default:
		err = fmt.Errorf("no such method")
	}
	return
}

// returns string
func MakeReqToDialogFlow(message string) (answer string, err error) {
	req, err := json.Marshal(obj{"message": message})
	if err != nil {
		fmt.Println("MakeReqToDialogFlow() -> json.Marshal error:", err.Error())
		return
	}

	resp, err := MakeHttpReq(dialogFlowerUrl+"/getAnswer", "POST", req)
	if err != nil {
		fmt.Println("MakeReqToDialogFlow() -> makeHttpReq(/getAnswer) error:", err.Error())
		return
	}

	var respStruct struct {
		Answer string `json:"answer"`
		Err    string `json:"err"`
	}

	if err := json.Unmarshal(resp, &respStruct); err != nil {
		fmt.Println("MakeReqToDialogFlow() -> unmarshal error:", err.Error())
		return "", err
	}

	if respStruct.Err != "" {
		fmt.Println("MakeReqToDialogFlow() -> got an error from dialogflow:", err.Error())
		return "", fmt.Errorf(respStruct.Err)
	}
	return respStruct.Answer, nil
}

// MakeReqToTost - makes req to tost service
func MakeReqToTost(method string, data []byte) (answer []byte, err error) {
	path := fmt.Sprintf("%s/%s", tostUrl, method)
	switch method {
	case "getRandomTost":
		answer, err = MakeHttpReq(path, "GET", data)
	default:
		err = fmt.Errorf("no such method")
	}
	return
}

// MakeHttpReq - func for sending http req with given path, method(get or post!) and data
func MakeHttpReq(path, method string, data []byte) (answer []byte, err error) {
	var resp *http.Response
	switch method {
	case "GET":
		resp, err = http.Get(path)
	case "POST":
		resp, err = http.Post(path, "application/json", bytes.NewReader(data))
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
