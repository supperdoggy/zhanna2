package main

import (
	"bytes"
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
