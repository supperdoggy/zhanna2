package tostdata

type GetRandomTostResp struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Err  string `json:"err"`
}

type GetTostByIdReq struct {
	ID int `json:"id" bson:"_id"`
}

type GetTostByIdResp struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Err  string `json:"err"`
}

type DeleteTostReq struct {
	ID int `json:"id" bson:"_id"`
}

type DeleteTostResp struct {
	Err string `json:"err"`
	OK  bool   `json:"ok"`
}

type AddTostReq struct {
	Text string `json:"text"`
}

type AddTostResp struct {
	Err string `json:"err"`
	OK  bool   `json:"ok"`
}
