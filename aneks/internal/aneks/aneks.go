package aneks

import (
	"github.com/supperdoggy/superSecretDevelopement/aneks/internal/db"
	aneksdata "github.com/supperdoggy/superSecretDevelopement/structs/request/aneks"
)

type AneksService struct {
	DB *db.DbStruct
}

func (s *AneksService) GetRandomAnek() (resp aneksdata.GetRandomAnekResp) {
	a, err := s.DB.GetRandomAnek()
	if err != nil {
		resp.Err = err.Error()
		return
	}

	resp.Text = a.Text
	resp.ID = a.Id
	return resp
}

func (s *AneksService) GetAnekByID(req aneksdata.GetAnekByIdReq) (resp aneksdata.GetAnekByIdResp) {
	result, err := s.DB.GetAnekById(req.ID)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Anek = result
	return
}

func (s *AneksService) DeleteAnekByID(req aneksdata.DeleteAnekByIDReq) (resp aneksdata.DeleteAnekByIDResp) {
	err := s.DB.DeleteAnek(req.ID)
	if err != nil {
		resp.Err = err.Error()
		return
	}

	resp.OK = true
	return
}

func (s *AneksService) AddAnek(req aneksdata.AddAnekReq) (resp aneksdata.AddAnekResp) {
	if req.Text == "" {
		resp.Err = "text field cant be empty"
		return
	}
	err := s.DB.AddAnek(req.Text)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.OK = true
	return
}
