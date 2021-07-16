package dialogflowdata

type GetAnswerReq struct {
	Text string `json:"text"`
	ID   string    `json:"id"`
}

type GetAnswerResp struct {
	Answer string `json:"answer"`
	Err    string `json:"err"`
}
