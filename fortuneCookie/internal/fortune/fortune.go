package fortune

import (
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	fortunedata "github.com/supperdoggy/superSecretDevelopement/structs/request/fortune"
	"go.uber.org/zap"
)

type (
	Service struct {
		db     db.IDbStruct
		logger *zap.Logger
	}
	IService interface {
		GetRandomFortuneCookie() (resp fortunedata.GetRandomFortuneCookieResp)
	}
)

func NewService(logger *zap.Logger, db db.IDbStruct) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (s Service) GetRandomFortuneCookie() (resp fortunedata.GetRandomFortuneCookieResp) {
	var cookie structs.Cookie

	cookie, err := s.db.GetRandomFortune()
	if err != nil {
		s.logger.Error("error GetRandomFortune", zap.Error(err), zap.Any("resp", resp))
		resp.Err = err.Error()
		return
	}
	resp.Text = cookie.Text
	resp.ID = cookie.ID
	return
}
