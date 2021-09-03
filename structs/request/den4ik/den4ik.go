package den4ikdata

import "github.com/supperdoggy/superSecretDevelopement/structs"

type GetCardReq struct {
	SessionID int `json:"session_id"`
}

type GetCardResp struct {
	Card         structs.Card `json:"card"`
	SessionIsNew bool         `json:"session_is_new"`
	SessionEnd   bool         `json:"session_end"`
	Err          string       `json:"err"`
}

type ResetSessionReq struct {
	SessionID int `json:"session_id"`
}

type ResetSessionResp struct {
	OK  bool   `json:"ok"`
	Err string `json:"err"`
}
