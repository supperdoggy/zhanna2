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

type obj map[string]interface{}
type arr []interface{}

type Service struct {
	DB *db.DbStruct
	Logger *zap.Logger
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

	if err := s.DB.AddFlower(flowerToAdd); err != nil {
		s.Logger.Error("error when AddFlower", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		return resp, err
	}
	resp.OK = true
	return
}

func (s *Service) RemoveFlower(req flowersdata.RemoveFlowerReq) (resp flowersdata.RemoveFlowerResp, err error) {
	err = s.DB.RemoveFlower(req.ID)
	if err != nil {
		s.Logger.Error("error when RemoveFlower", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		return
	}
	resp.OK = true
	return
}

//goland:noinspection GoNilness
func (s *Service) GrowFlower(req flowersdata.GrowFlowerReq) (resp flowersdata.GrowFlowerResp, err error) {
	flower, err := s.DB.GetUserCurrentFlower(req.ID)
	if err != nil && err.Error() != "not found" {
		s.Logger.Error("error getting flower GetUserCurrentFlower", zap.Error(err), zap.Any("req", req))
		resp.Err = "error getting flower"
		return
	} else if err != nil && err.Error() == "not found" { // not found flower, creating new
		s.Logger.Info("creating new flower to grow", zap.Any("req", req))
		flower, err = s.DB.GetRandomFlower()
		if err != nil {
			s.Logger.Error("error getting random flower", zap.Error(err), zap.Any("req", req))
			resp.Err = err.Error()
			return
		}
		flower.Owner = req.ID
		flower.Grew = uint8(rand.Intn(31)) + 1 // so its not possible to get 0
		flower.HP += flower.Grew
	}

	// check if flower died
	if !req.NonDying {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(101)
		died := num >= 0 && num <= cfg.DyingChance
		if died {
			s.Logger.Info("flower died", zap.Any("req", req), zap.Any("flower", flower))
			resp.Err = "flower died"
			flower.Dead = true
			resp.Flower = flower
			if err := s.DB.EditUserFlower(flower); err != nil {
				s.Logger.Error("error EditUserFlower",
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
	flower.Grew = uint8(rand.Intn(31) + extraGrow)
	flower.HP += flower.Grew

	if flower.HP > 100 {
		flower.HP = 100
	}
	flower.LastTimeGrow = time.Now()

	if err := s.DB.EditUserFlower(flower); err != nil {
		s.Logger.Error("error EditUserFlower",
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
	flowers, err := s.DB.GetAllUserFlowers(req.ID)
	if err != nil {
		s.Logger.Error("error GetAllUserFlowers", zap.Error(err), zap.Any("req", req))
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
			Name   string `json:"name"`
			Amount int    `json:"amount"`
		}{Name: v.Icon + " " + v.Name, Amount: count[v.Name+v.Icon]})
		types[v.Name+v.Icon] = false
	}

	var total int
	for _, v := range count {
		total += v
	}
	var last uint8
	flower, err := s.DB.GetUserCurrentFlower(req.ID)
	if err != nil {
		s.Logger.Error("error when GetUserCurrentFlower", zap.Error(err), zap.Any("req", req))
	}
	last = flower.HP
	resp.Total = total
	resp.Last = last
	return
}

func (s *Service) CanGrowFlower(req flowersdata.CanGrowFlowerReq) (resp flowersdata.CanGrowFlowerResp, err error) {
	flower, err := s.DB.GetUserCurrentFlower(req.ID)
	if err != nil {
		// if we cant find flower in the collection we return true
		if err.Error() == "not found" {
			resp.Answer = true
			err = nil
			return
		}
		s.Logger.Error("error GetUserCurrentFlower", zap.Error(err), zap.Any("req", req))
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
		err := s.DB.RemoveUserFlower(defaultCfg.Obj{"owner": req.ID, "hp": obj{"$ne": 100}})
		if err != nil {
			s.Logger.Error("error RemoveUserFlower", zap.Error(err), zap.Any("req", req))
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
	flowersCount, err := s.DB.CountFlowers(req.ID)
	if err != nil {
		s.Logger.Error("error CountFlowers", zap.Error(err), zap.Any("req", req))
		resp.Err = "error getting flowers"
		return
	}
	resp.Total = flowersCount
	return
}

func (s *Service) GetLastFlower(req flowersdata.GetLastFlowerReq) (resp flowersdata.GetLastFlowerResp, err error) {
	flower, err := s.DB.GetUserCurrentFlower(req.ID)
	if err != nil {
		s.Logger.Error("error GetUserCurrentFlower", zap.Error(err), zap.Any("req", req))
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
	result, err := s.DB.UserFlowerSlice(req.ID)
	if err != nil {
		s.Logger.Error("error UserFlowerSlice", zap.Error(err), zap.Any("req", req))
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

// TODO: simplify
func (s *Service) GiveFlower(req flowersdata.GiveFlowerReq) (resp flowersdata.GiveFlowerResp, err error) {
	if req.Owner == 0 || req.Reciever == 0 {
		resp.Err = "empty id"
		return resp, errors.New("empty id")
	}

	var f structs.Flower
	if req.Last {
		// getting flowers
		flowers, err := s.DB.GetAllUserFlowers(req.Owner)
		if err != nil { // if has no flower
			s.Logger.Error("error GetAllUserFlowers", zap.Error(err), zap.Any("req", req))

			resp.Err = "user has no flowers"
			return resp, errors.New("user has no flowers")
		}
		rand.Seed(time.Now().UnixNano())
		if len(flowers) != 0 {
			f = flowers[len(flowers)-1]
		}
	} else {
		f, _ = s.DB.GetUserFlowerById(req.ID)
	}
	if f.ID == 0 {
		resp.Err = "user has no flowers"
		return resp, errors.New("user has no flowers")
	}

	if req.Reciever == cfg.TestId || req.Reciever == cfg.ProdId {
		f.Owner = cfg.ZhannaId
	} else {
		f.Owner = req.Reciever
	}

	if err := s.DB.EditUserFlower(f); err != nil {
		s.Logger.Error("error EditUserFlower", zap.Error(err), zap.Any("req", req), zap.Any("flower", f))

		resp.Err = err.Error()
		return resp, err
	}
	resp.Flower = f
	return
}

func (s *Service) GetFlowerTypes() (resp flowersdata.GetFlowerTypesResp, err error) {
	flowers, err := s.DB.GetAllFlowers()
	if err != nil {
		s.Logger.Error("error GetAllFlowers", zap.Error(err))

		resp.Err = err.Error()
		return
	}
	resp.Result = flowers
	return
}
