package main

import "log"

func checkAdmin(id int) (bool, error) {
	data, err := MakeAdminHTTPReq("isAdmin", obj{"id": id})
	if err != nil {
		log.Println("admin_auth.go -> checkAdmin() -> isAdmin method req error:", err)
		return false, err
	}
	var resp struct {
		Err    string `json:"err"`
		Result string `json:"result"`
	}
	if resp.Err != "" {
		log.Println("admin_auth.go -> checkAdmin() -> resp error:", resp.Err)
		return false, resp.Err
	}
	return resp.Result, ""
}
