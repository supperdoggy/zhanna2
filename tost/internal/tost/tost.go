package tost

import (
	tostdata "github.com/supperdoggy/superSecretDevelopement/structs/request/tost"
	"github.com/supperdoggy/superSecretDevelopement/tost/internal/db"
)

type TostService struct {
	DB *db.DbStruct
}

func (s TostService) GetRandomTost() (resp tostdata.GetRandomTostResp) {
	a, err := s.DB.GetRandomTost()
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Text = a.Text
	resp.ID = a.ID
	return
}

func (s TostService) GetTostById(req tostdata.GetTostByIdReq) (resp tostdata.GetTostByIdResp) {
	result, err := s.DB.GetTostById(req.ID)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Text = result.Text
	resp.ID = result.ID
	return
}

func (s TostService) DeleteTost(req tostdata.DeleteTostReq) (resp tostdata.DeleteTostResp) {
	err := s.DB.DeleteTost(req.ID)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.OK = true
	return
}

func (s TostService) AddTost(req tostdata.AddTostReq) (resp tostdata.AddTostResp) {
	if req.Text == "" {
		resp.Err = "text field is empty"
		return
	}

	err := s.DB.AddTost(req.Text)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.OK = true
	return
}
