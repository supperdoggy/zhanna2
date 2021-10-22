package flowersdata

import "github.com/supperdoggy/superSecretDevelopement/structs"

type AddNewFlowerReq struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Type string `json:"type"`
}

type AddNewFlowerResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}

type RemoveFlowerReq struct {
	ID uint64 `json:"id"`
}

type RemoveFlowerResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}

type GrowFlowerReq struct {
	ID       int  `json:"id"`
	NonDying bool `json:"nonDying"`
	MsgCount int  `json:"msg_count"`
}

type GrowFlowerResp struct {
	Flower structs.Flower `json:"flower"`
	Extra  int            `json:"extra"`
	Err    string         `json:"err"`
}

type GetUserFlowersReq struct {
	ID int `json:"id" bson:"owner"`
}

type GetUserFlowersResp struct {
	Flowers []struct {
		NameAndIcon string `json:"name_and_icon"`
		Name        string `json:"name"`
		Amount      int    `json:"amount"`
	} `json:"flowers"`
	Total int    `json:"total"`
	Last  uint8  `json:"last"`
	Err   string `json:"err"`
}

type CanGrowFlowerReq struct {
	ID int `json:"id" bson:"owner"`
}

type CanGrowFlowerResp struct {
	Answer bool   `json:"answer"`
	Err    string `json:"err"`
}

type RemoveUserFlowerReq struct {
	ID      int  `json:"id" bson:"owner"`
	Current bool `json:"current"`
}

type RemoveUserFlowerResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}

type GetUserFlowerTotalReq struct {
	ID int `json:"id" bson:"owner"`
}

type GetUserFlowerTotalResp struct {
	Total int    `json:"total"`
	Err   string `json:"err"`
}

type GetLastFlowerReq struct {
	ID int `json:"id" bson:"owner"`
}

type GetLastFlowerResp struct {
	Flower structs.Flower `json:"flower"`
	Err    string         `json:"err"`
}

type UserFlowerSliceReq struct {
	ID []int `json:"id" bson:"owner"`
}

type UserFlowerSliceResp struct {
	Result []struct {
		Key   int `json:"id"`
		Value int `json:"total"`
	} `json:"result"`

	Err string `json:"err"`
}

type GiveFlowerReq struct {
	Owner    int    `json:"owner"`
	Reciever int    `json:"reciever"`
	Last     bool   `json:"last"`
	ID       string `json:"id"`
}

type GiveFlowerResp struct {
	Err    string         `json:"err"`
	Flower structs.Flower `json:"flower"`
}

type GetFlowerTypesResp struct {
	Result []structs.Flower `json:"result"`
	Err    string           `json:"err"`
}
