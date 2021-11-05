package service

import (
	"errors"
	"github.com/supperdoggy/superSecretDevelopement/flowers/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	flowersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/flowers"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"go.uber.org/zap"
	"math"
	"math/rand"
	"sort"
	"time"
)

type (
	Service struct {
		db     db.IDbStruct
		logger *zap.Logger
	}
	IService interface {
		AddNewFlower(req flowersdata.AddNewFlowerReq) (resp flowersdata.AddNewFlowerResp, err error)
		RemoveFlower(req flowersdata.RemoveFlowerReq) (resp flowersdata.RemoveFlowerResp, err error)
		GrowFlower(req flowersdata.GrowFlowerReq) (resp flowersdata.GrowFlowerResp, err error)
		GetUserFlowers(req flowersdata.GetUserFlowersReq) (resp flowersdata.GetUserFlowersResp, err error)
		CanGrowFlower(req flowersdata.CanGrowFlowerReq) (resp flowersdata.CanGrowFlowerResp, err error)
		RemoveUserFlower(req flowersdata.RemoveUserFlowerReq) (resp flowersdata.RemoveUserFlowerResp, err error)
		GetUserFlowerTotal(req flowersdata.GetUserFlowerTotalReq) (resp flowersdata.GetUserFlowerTotalResp, err error)
		GetLastFlower(req flowersdata.GetLastFlowerReq) (resp flowersdata.GetLastFlowerResp, err error)
		UserFlowerSlice(req flowersdata.UserFlowerSliceReq) (resp flowersdata.UserFlowerSliceResp, err error)
		GiveFlower(req flowersdata.GiveFlowerReq) (resp flowersdata.GiveFlowerResp, err error)
		GetFlowerTypes() (resp flowersdata.GetFlowerTypesResp, err error)
	}
)

func NewService(db db.IDbStruct, l *zap.Logger) *Service {
	return &Service{
		db:     db,
		logger: l,
	}
}

func (s *Service) AddNewFlower(req flowersdata.AddNewFlowerReq) (resp flowersdata.AddNewFlowerResp, err error) {
	if req.Name == "" || req.Icon == "" || req.Type == "" {
		resp.Err = "fill all fields"
		return resp, errors.New("fill all fields")
	}

	flowerToAdd := structs.Flower{
		Name:         req.Name,
		Icon:         req.Icon,
		Type:         req.Type,
		CreationTime: time.Now(),
	}

	if err := s.db.AddFlower(flowerToAdd); err != nil {
		s.logger.Error("error when AddFlower", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		return resp, err
	}
	resp.OK = true
	return
}

func (s *Service) RemoveFlower(req flowersdata.RemoveFlowerReq) (resp flowersdata.RemoveFlowerResp, err error) {
	err = s.db.RemoveFlower(req.ID)
	if err != nil {
		s.logger.Error("error when RemoveFlower", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		return
	}
	resp.OK = true
	return
}

//goland:noinspection GoNilness
func (s *Service) GrowFlower(req flowersdata.GrowFlowerReq) (resp flowersdata.GrowFlowerResp, err error) {
	flower, err := s.db.GetUserCurrentFlower(req.ID)
	var flowerIsNew bool = false
	if err != nil && err.Error() != "not found" {
		s.logger.Error("error getting flower GetUserCurrentFlower", zap.Error(err), zap.Any("req", req))
		resp.Err = "error getting flower"
		return
	} else if err != nil && err.Error() == "not found" { // not found flower, creating new
		s.logger.Info("creating new flower to grow", zap.Any("req", req))
		flower, err = s.db.GetRandomFlower()
		if err != nil {
			s.logger.Error("error getting random flower", zap.Error(err), zap.Any("req", req))
			resp.Err = err.Error()
			return
		}
		flower.Owner = req.ID
		flower.CreationTime = time.Now()
		flowerIsNew = true
	}

	// check if flower died
	if !req.NonDying {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(101)
		died := num >= 0 && num <= cfg.DyingChance
		if died && !flowerIsNew {
			s.logger.Info("flower died", zap.Any("req", req), zap.Any("flower", flower))
			resp.Err = "flower died"
			flower.Dead = true
			resp.Flower = flower
			if err := s.db.EditUserFlower(flower); err != nil {
				s.logger.Error("error EditUserFlower",
					zap.Error(err),
					zap.Any("req", req),
					zap.Any("flower", flower),
				)

				resp.Err = err.Error()
				return resp, err
			}
			return
		}
	}
	// if did not die

	// add extra grow output for user
	extraGrow := int(math.Round(float64(req.MsgCount) * cfg.Message_multiplyer))
	if extraGrow > 20 {
		extraGrow = 20
	}

	grew := rand.Intn(31) + 1
	flower.Grew = uint8(grew + extraGrow)
	flower.HP += uint8(grew + extraGrow)

	if flower.HP > 100 {
		flower.HP = 100
	}
	flower.LastTimeGrow = time.Now()

	err = nil
	if flowerIsNew {
		err = s.db.CreateUserFlower(flower)
	} else {
		err = s.db.EditUserFlower(flower)
	}

	if err != nil {
		s.logger.Error("error grow flower",
			zap.Error(err),
			zap.Any("req", req),
			zap.Any("flower", flower),
		)

		resp.Err = err.Error()
		return resp, err
	}

	resp.Flower = flower
	resp.Extra = extraGrow
	return
}

func (s *Service) GetUserFlowers(req flowersdata.GetUserFlowersReq) (resp flowersdata.GetUserFlowersResp, err error) {
	flowers, err := s.db.GetAllUserFlowers(req.ID)
	if err != nil {
		s.logger.Error("error GetAllUserFlowers", zap.Error(err), zap.Any("req", req))
		resp.Err = "error getting flowers"
		return
	}

	count := map[string]int{}
	// only to check what flowers we have
	types := map[string]bool{}

	for _, v := range flowers {
		count[v.Name+v.Icon]++
		types[v.Name+v.Icon] = true
	}

	for _, v := range flowers {
		if find := types[v.Name+v.Icon]; !find {
			continue
		}
		resp.Flowers = append(resp.Flowers, struct {
			NameAndIcon string `json:"name_and_icon"`
			Name        string `json:"name"`
			Amount      int    `json:"amount"`
		}{Name: v.Name, NameAndIcon: v.Icon + " " + v.Name, Amount: count[v.Name+v.Icon]})
		types[v.Name+v.Icon] = false
	}

	var total int
	for _, v := range count {
		total += v
	}
	var last uint8
	flower, err := s.db.GetUserCurrentFlower(req.ID)
	if err != nil {
		s.logger.Error("error when GetUserCurrentFlower", zap.Error(err), zap.Any("req", req))
	}
	last = flower.HP
	resp.Total = total
	resp.Last = last
	return
}

func (s *Service) CanGrowFlower(req flowersdata.CanGrowFlowerReq) (resp flowersdata.CanGrowFlowerResp, err error) {
	flower, err := s.db.GetUserCurrentFlower(req.ID)
	if err != nil {
		// if we cant find flower in the collection we return true
		if err.Error() == "not found" {
			resp.Answer = true
			err = nil
			return
		}
		s.logger.Error("error GetUserCurrentFlower", zap.Error(err), zap.Any("req", req))
		// if we cant find due to mongo error then return error
		resp.Err = "got flower error"
		return
	}
	// if passed GrowTimeout hours
	canGrow := int(time.Now().Sub(flower.LastTimeGrow).Hours())/cfg.GrowTimeout >= 1
	resp.Answer = canGrow
	return
}

func (s *Service) RemoveUserFlower(req flowersdata.RemoveUserFlowerReq) (resp flowersdata.RemoveUserFlowerResp, err error) {
	if req.ID == 0 {
		resp.Err = "no id field"
		return resp, errors.New("no id field")
	}
	if req.Current {
		err := s.db.RemoveUserFlower(defaultCfg.Obj{"owner": req.ID, "hp": defaultCfg.Obj{"$ne": 100}})
		if err != nil {
			s.logger.Error("error RemoveUserFlower", zap.Error(err), zap.Any("req", req))
			resp.Err = "error removing"
			return resp, err
		}
		resp.OK = true
		return resp, err
	}
	// todo: remove random flower
	resp.Err = "not implemented random"
	return resp, errors.New(resp.Err)
}

func (s *Service) GetUserFlowerTotal(req flowersdata.GetUserFlowerTotalReq) (resp flowersdata.GetUserFlowerTotalResp, err error) {
	flowersCount, err := s.db.CountFlowers(req.ID)
	if err != nil {
		s.logger.Error("error CountFlowers", zap.Error(err), zap.Any("req", req))
		resp.Err = "error getting flowers"
		return
	}
	resp.Total = flowersCount
	return
}

func (s *Service) GetLastFlower(req flowersdata.GetLastFlowerReq) (resp flowersdata.GetLastFlowerResp, err error) {
	flower, err := s.db.GetUserCurrentFlower(req.ID)
	if err != nil {
		s.logger.Error("error GetUserCurrentFlower", zap.Error(err), zap.Any("req", req))
		resp.Err = "error getting flowers"
		return
	}
	resp.Flower = flower
	return
}

func (s *Service) UserFlowerSlice(req flowersdata.UserFlowerSliceReq) (resp flowersdata.UserFlowerSliceResp, err error) {
	if len(req.ID) == 0 {
		resp.Err = "empty id slice"
		return resp, errors.New(resp.Err)
	}
	result, err := s.db.UserFlowerSlice(req.ID)
	if err != nil {
		s.logger.Error("error UserFlowerSlice", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		return resp, err
	}

	var scores = make(map[int]int)
	for _, v := range result {
		scores[v.Owner]++
	}
	for k, v := range scores {
		resp.Result = append(resp.Result, struct {
			Key   int `json:"id"`
			Value int `json:"total"`
		}{k, v})
	}

	sort.Slice(resp.Result, func(i, j int) bool {
		return resp.Result[i].Value > resp.Result[j].Value
	})
	return
}

func (s *Service) GiveFlower(req flowersdata.GiveFlowerReq) (resp flowersdata.GiveFlowerResp, err error) {
	if req.Owner == 0 || req.Reciever == 0 {
		resp.Err = "empty id"
		return resp, errors.New("empty id")
	}

	var f structs.Flower
	if req.Last {
		// getting flowers
		f, err = s.db.GetLastUserFlower(req.Owner)
	} else {
		f, err = s.db.GetUserFlowerByName(req.Owner, req.ID)
	}

	if f.ID == 0 || err != nil {
		resp.Err = err.Error()
		return resp, errors.New(resp.Err)
	}

	if req.Reciever == cfg.TestId || req.Reciever == cfg.ProdId {
		f.Owner = cfg.ZhannaId
	} else {
		f.Owner = req.Reciever
	}

	if err := s.db.EditUserFlower(f); err != nil {
		s.logger.Error("error EditUserFlower", zap.Error(err), zap.Any("req", req), zap.Any("flower", f))

		resp.Err = err.Error()
		return resp, err
	}
	resp.Flower = f
	return
}

func (s *Service) GetFlowerTypes() (resp flowersdata.GetFlowerTypesResp, err error) {
	flowers, err := s.db.GetAllFlowers()
	if err != nil {
		s.logger.Error("error GetAllFlowers", zap.Error(err))

		resp.Err = err.Error()
		return
	}
	resp.Result = flowers
	return
}
