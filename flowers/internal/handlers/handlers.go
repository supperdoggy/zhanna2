package handlers

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/flowers/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	flowersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/flowers"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"time"

	ai "github.com/night-codes/mgo-ai"

	"github.com/gin-gonic/gin"
)

type obj map[string]interface{}

type Handlers struct {
	DB *db.DbStruct
}

// adds new flower type
func (h Handlers) AddNewFlower(c *gin.Context) {
	var req flowersdata.AddNewFlowerReq
	var resp flowersdata.AddNewFlowerResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> addNewFlower() -> binding error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if req.Name == "" || req.Icon == "" || req.Type == "" {
		resp.Err = "fill all fields"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	flowerToAdd := structs.Flower{
		ID:           ai.Next(h.DB.FlowerCollection.Name),
		Name:         req.Name,
		Icon:         req.Icon,
		Type:         req.Type,
		CreationTime: time.Now(),
	}

	if err := h.DB.AddFlower(flowerToAdd); err != nil {
		fmt.Println("handlers.go -> addNewFlower() -> addFlower(req) error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.OK = true
	c.JSON(http.StatusOK, resp)
}

// removes flower type
func (h Handlers) RemoveFlower(c *gin.Context) {
	var req flowersdata.RemoveFlowerReq
	var resp flowersdata.RemoveFlowerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println("handlers.go -> removeFlower() -> bind error:", err.Error(), string(data))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := h.DB.RemoveFlower(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> removeFlower() -> removeFlower() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.OK = true
	c.JSON(http.StatusOK, resp)
}

// grows user flower
func (h Handlers) GrowFlower(c *gin.Context) {
	var req flowersdata.GrowFlowerReq
	var resp flowersdata.GrowFlowerResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> growFlowerReq) -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	flower, err := h.DB.GetUserFlower(req.ID)
	if err != nil && err.Error() != "not found" {
		log.Println("error getting")
		fmt.Println("handlers.go -> growFlowerReq) -> getUserFlower() error:", err.Error())
		resp.Err = "error getting flower"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	// not found flower, creating new
	if err != nil && err.Error() == "not found" {
		log.Println("creating new")
		flower, err = h.DB.GetRandomFlower()
		if err != nil {
			fmt.Println("handlers.go -> growFlowerReq) -> getRandomFlower() error:", err.Error())
			resp.Err = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		flower.ID = ai.Next(h.DB.UserFlowerDataCollection.Name)
		flower.Owner = req.ID
		flower.Grew = uint8(rand.Intn(31)) + 1 // so its not possible to get 0
		flower.HP += flower.Grew
	}

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

	if _, err := h.DB.UserFlowerDataCollection.Upsert(obj{"_id": flower.ID}, flower); err != nil {
		fmt.Println("handlers.go -> growFlowerReq) -> Upsert() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Flower = flower
	resp.Extra = extraGrow
	c.JSON(http.StatusOK, resp)

}

// returns map of user flowers and quantity of different type
func (h Handlers) GetUserFlowers(c *gin.Context) {
	var req flowersdata.GetUserFlowersReq
	var resp flowersdata.GetUserFlowersResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getUserFlowers() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	flowers, err := h.DB.GetAllUserFlowersMap(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getUserFlowers() -> getAllUserFlowers() error:", err.Error())
		resp.Err = "error getting flowers"
		c.JSON(http.StatusBadRequest, resp)
	}

	var total int
	for _, v := range flowers {
		total += v
	}
	var last uint8
	flower, err := h.DB.GetUserFlower(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getUserFlowers() -> getUserFlower() error:", err.Error())
	}
	last = flower.HP
	resp.Flowers = flowers
	resp.Total = total
	resp.Last = last
	c.JSON(http.StatusOK, resp)
}

// returns bool value if user can grow flower
func (h Handlers) CanGrowFlower(c *gin.Context) {
	var req flowersdata.CanGrowFlowerReq
	var resp flowersdata.CanGrowFlowerResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> canGrowFlower() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	flower, err := h.DB.GetUserFlower(req.ID)
	if err != nil {
		// if we cant find flower in the collection we return true
		if err.Error() == "not found" {
			resp.Answer = true
			c.JSON(http.StatusOK, resp)
			return
		}
		// if we cant find due to mongo error then return error
		fmt.Println("handlers.go -> canGrowFlower() -> getUserFlower() error:", err.Error())
		resp.Err = "got flower error"
		c.JSON(http.StatusBadRequest, resp)
	}
	// if passed GrowTimeout hours
	canGrow := int(time.Now().Sub(flower.LastTimeGrow).Hours())/cfg.GrowTimeout >= 1
	resp.Answer = canGrow
	c.JSON(http.StatusOK, resp)
}

// removeUserFlower - removes current user flower
func (h Handlers) RemoveUserFlower(c *gin.Context) {
	var req flowersdata.RemoveUserFlowerReq
	var resp flowersdata.RemoveUserFlowerResp

	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> removeUserFlower() -> c.Bind() error:", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if req.ID == 0 {
		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if req.Current {
		err := h.DB.UserFlowerDataCollection.Remove(obj{"owner": req.ID, "hp": obj{"$ne": 100}})
		if err != nil {
			fmt.Println("handlers.go -> removeUserFlower() -> Remove() error:", err.Error())
			resp.Err = "error removing"
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		resp.OK = true
		c.JSON(http.StatusOK, resp)
		return
	}
	// todo: remove random flower
}

// returns int quantity of user grown flowers
func (h Handlers) GetUserFlowerTotal(c *gin.Context) {
	var req flowersdata.GetUserFlowerTotalReq
	var resp flowersdata.GetUserFlowerTotalResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getUserFlowerTotal() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	flowersCount, err := h.DB.CountFlowers(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getUserFlowerTotal() -> getAllUserFlowers() error:", err.Error())
		resp.Err = "error getting flowers"
		c.JSON(http.StatusBadRequest, resp)
	}
	resp.Total = flowersCount
	c.JSON(http.StatusOK, resp)
}

func (h Handlers) GetLastFlower(c *gin.Context) {
	var req flowersdata.GetLastFlowerReq
	var resp flowersdata.GetLastFlowerResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getUserFlowerTotal() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	flower, err := h.DB.GetUserFlower(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getLastFlower() -> getUserFlower() error:", err.Error())
		resp.Err = "error getting flowers"
		c.JSON(http.StatusBadRequest, resp)
	}
	resp.Flower = flower
	c.JSON(http.StatusOK, resp)
}

// userFlowerSlice - returns slice of users flowers
func (h Handlers) UserFlowerSlice(c *gin.Context) {
	var req flowersdata.UserFlowerSliceReq
	var resp flowersdata.UserFlowerSliceResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getUserFlowerTotal() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if len(req.ID) == 0 {
		resp.Err = "empty id slice"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	// building query for request to mongo
	query := []obj{}
	for _, v := range req.ID {
		query = append(query, obj{"owner": v})
	}

	var result []structs.Flower
	// TODO: BUG: returns dead flowers
	if err := h.DB.UserFlowerDataCollection.Find(obj{"$or": query}).Select(obj{"owner": 1, "hp": 1}).All(&result); err != nil {
		fmt.Println("handlers.go -> userFlowerSlice() -> flower find error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var scores map[int]int = make(map[int]int)
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
	c.JSON(http.StatusOK, resp)
}

// gives flower to other user
func (h Handlers) GiveFlower(c *gin.Context) {
	var req flowersdata.GiveFlowerReq
	var resp flowersdata.GiveFlowerResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> giveRandomFlower() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if req.Owner == 0 || req.Reciever == 0 {
		resp.Err = "empty id"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var f structs.Flower
	if req.Last {
		fmt.Println(req)
		// getting flowers
		flowers, err := h.DB.GetAllUserFlowers(req.Owner)
		if err != nil { // if has no flower
			resp.Err = "user has no flowers"
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		rand.Seed(time.Now().UnixNano())
		if len(flowers) != 0 {
			f = flowers[len(flowers)-1]
		}
	} else {
		f, _ = h.DB.GetUserFlowerById(req.ID)
	}
	if f.ID == 0 {
		resp.Err = "user has no flowers"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if req.Reciever == cfg.TestId || req.Reciever == cfg.ProdId {
		f.Owner = cfg.ZhannaId
	} else {
		f.Owner = req.Reciever
	}

	if err := h.DB.EditUserFlower(f.ID, f); err != nil {
		fmt.Println("handlers.go -> giveRandomFlower() -> editFlower() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// getFlowerTypes - for admin, returns all flower types
func (h Handlers) GetFlowerTypes(c *gin.Context) {
	var resp flowersdata.GetFlowerTypesResp
	flowers, err := h.DB.GetAllFlowers()
	if err != nil {
		log.Println("handlers.go -> getFlowerTypes() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Result = flowers
	c.JSON(http.StatusOK, resp)
}
