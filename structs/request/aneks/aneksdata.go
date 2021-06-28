package aneksdata

import "github.com/supperdoggy/superSecretDevelopement/structs"

type GetAnekByIdReq struct {
	ID int `json:"id"`
}

type GetAnekByIdResp struct {
	Err  string       `json:"err"`
	Anek structs.Anek `json:"anek"`
}

type GetRandomAnekResp struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Err  string `json:"err"`
}

type DeleteAnekByIDReq struct {
	ID int `json:"id"`
}

type DeleteAnekByIDResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}

type AddAnekReq struct {
	Text string `json:"text"`
}

type AddAnekResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}
