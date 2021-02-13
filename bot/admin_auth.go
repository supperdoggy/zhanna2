package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func checkAdmin(id int) (bool, error) {
	if id == NeMoksID {
		return true, nil
	}
	data, err := MakeAdminHTTPReq("isAdmin", obj{"id": id})
	if err != nil {
		log.Println("admin_auth.go -> checkAdmin() -> isAdmin method req error:", err)
		return false, err
	}
	var resp struct {
		Err    string `json:"err"`
		Result bool   `json:"result"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return false, err
	}

	if resp.Err != "" {
		log.Println("admin_auth.go -> checkAdmin() -> resp error:", resp.Err)
		return false, fmt.Errorf(resp.Err)
	}
	return resp.Result, nil
}
