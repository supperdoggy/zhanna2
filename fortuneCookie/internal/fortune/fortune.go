package fortune

import (
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	fortunedata "github.com/supperdoggy/superSecretDevelopement/structs/request/fortune"
)

type Service struct {
	DB db.DbStruct
}

func (s Service) GetRandomFortuneCookie() (resp fortunedata.GetRandomFortuneCookieResp) {
	var cookie structs.Cookie

	cookie, err := s.DB.GetRandomFortune()
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Text = cookie.Text
	resp.ID = cookie.ID
	return
}
