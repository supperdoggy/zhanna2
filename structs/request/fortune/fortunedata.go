package fortunedata

type GetRandomFortuneCookieResp struct {
	Text string `json:"text"`
	ID   int32  `json:"id"`
	Err string `json:"err"`
}
