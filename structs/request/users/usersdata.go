package usersdata

import "github.com/supperdoggy/superSecretDevelopement/structs"

// admin handlers requests

type IsAdminReq struct {
	ID int `json:"id"`
}

type IsAdminResp struct {
	Result bool   `json:"result"`
	Err    string `json:"err"`
}

type AdminReq struct {
	ID int `json:"id"`
}

type AdminResp struct {
	Admin bool   `json:"admin"`
	OK    bool   `json:"ok"`
	Err   string `json:"err"`
}

type GetAllFlowerTypesResp struct {
	Result []structs.Flower `json:"result"`
	Err    string           `json:"err"`
}

type RemoveFlowerReq struct {
	ID uint64 `json:"id"`
}

type RemoveFlowerResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}

// user handlers requests

// todo think of a good solution))
type AddOrUpdateUserReq struct {
}

type AddOrUpdateUserResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}

type GetFortuneReq struct {
	ID int `json:"id" bson:"id" form:"id"`
}

type GetFortuneResp struct {
	Fortune structs.Cookie `json:"fortune"`
	Err     string         `json:"err"`
}

type GetRandomAnekReq struct {
	ID int `json:"id" bson:"id"`
}

type GetRandomAnekResp struct {
	structs.Anek
	Err string `json:"err"`
}

type GetRandomTostReq struct {
	ID int `json:"id" bson:"_id"`
}

type GetRandomTostResp struct {
	structs.Tost
	Err string `json:"err"`
}

type AddFlowerReq struct {
	Icon string `json:"icon" bson:"icon"`
	Name string `json:"name" bson:"name"`
	Type string `json:"type" bson:"type"`
}

type AddFlowerResp struct {
	// response flower type
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}

type FlowerReq struct {
	ID       int  `json:"id"`
	NonDying bool `json:"nonDying"`
	MsgCount int  `json:"msg_count"`
}

type FlowerResp struct {
	structs.Flower
	Up    uint8  `json:"up"`
	Grew  bool   `json:"grew"`
	Extra int    `json:"extra"`
	Err   string `json:"err"`
}

type DialogFlowReq struct {
	Text string `json:"text"`
	ID   string    `json:"id"`
}

type DialogFlowResp struct {
	Answer string `json:"answer"`
	Err    string `json:"err"`
}

type MyFlowersReq struct {
	ID int `json:"id"`
}

type MyFlowersResp struct {
	Flowers map[string]int `json:"flowers"`
	Last    uint8          `json:"last"`
	Total   int            `json:"total"`
	Err     string         `json:"err"`
}

type GiveFlowerReq struct {
	ID       uint64 `json:"id"`
	Owner    int    `json:"owner"`
	Reciever int    `json:"reciever"`
	// this is check for last or number of lasts
	Last  bool `json:"last"`
	Count int  `json:"count"`
}

type GiveFlowerResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}

type FlowertopReq struct {
	ChatId int `json:"chatid"`
}

type FlowertopResp struct {
	Result []struct {
		Username string `json:"username"`
		Total    int    `json:"total"`
	} `json:"result"`
	Err string `json:"err"`
}

type GetRandomNHIEreq struct {
	ID int `json:"id"`
}

type GetRandomNHIEresp struct {
	Err    string       `json:"err"`
	Result structs.NHIE `json:"result"`
}
