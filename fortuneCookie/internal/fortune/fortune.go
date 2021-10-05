package fortune

import (
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	fortunedata "github.com/supperdoggy/superSecretDevelopement/structs/request/fortune"
	"go.uber.org/zap"
)

type Service struct {
	DB     db.DbStruct
	Logger *zap.Logger
}

func (s Service) GetRandomFortuneCookie() (resp fortunedata.GetRandomFortuneCookieResp) {
	var cookie structs.Cookie

	cookie, err := s.DB.GetRandomFortune()
	if err != nil {
		s.Logger.Error("error GetRandomFortune", zap.Error(err), zap.Any("resp", resp))
		resp.Err = err.Error()
		return
	}
	resp.Text = cookie.Text
	resp.ID = cookie.ID
	return
}
