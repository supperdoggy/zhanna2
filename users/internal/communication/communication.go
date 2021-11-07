package communication

import (
	"bytes"
	"encoding/json"
	"fmt"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	anekscfg "github.com/supperdoggy/superSecretDevelopement/structs/services/aneks"
	flowerscfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	tostcfg "github.com/supperdoggy/superSecretDevelopement/structs/services/tost"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type obj map[string]interface{}

// MakeReqToAnek - makes req to anek service
func MakeReqToAnek(method string, req, resp interface{}) (err error) {
	path := cfg.AnekURL + method
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	var answer []byte
	switch method {
	case anekscfg.GetRandomAnekURL:
		answer, err = MakeHttpReq(path, "GET", data)
	default:
		err = fmt.Errorf("no such method")
	}
	if err != nil {
		return
	}
	return json.Unmarshal(answer, resp)
}

// MakeReqToFlowers - makes req to flowers service
func MakeReqToFlowers(method string, req, resp interface{}) (err error) {
	path := cfg.FlowersURL + method
	reqData, err := json.Marshal(req)
	if err != nil {
		return
	}
	var answer []byte
	switch method {
	case flowerscfg.AddNewFlowerURL:
		answer, err = MakeHttpReq(path, "POST", reqData)
	case flowerscfg.GrowFlowerURL:
		answer, err = MakeHttpReq(path, "POST", reqData)
	case flowerscfg.CanGrowFlowerURL:
		answer, err = MakeHttpReq(path, "POST", reqData)
	case flowerscfg.GetUserFlowersURL:
		answer, err = MakeHttpReq(path, "POST", reqData)
	case flowerscfg.GiveFlowerURL:
		answer, err = MakeHttpReq(path, "POST", reqData)
	case flowerscfg.UserFlowerSliceURL:
		answer, err = MakeHttpReq(path, "POST", reqData)
	case flowerscfg.GetFlowerTypesURL:
		answer, err = MakeHttpReq(path, "GET", nil)
	case flowerscfg.RemoveFlowerURL:
		answer, err = MakeHttpReq(path, "POST", reqData)
	case flowerscfg.AddUserFlowerURL:
		answer, err = MakeHttpReq(path, "POST", reqData)
	default:
		err = fmt.Errorf("no such method")
	}
	if err != nil {
		return
	}
	return json.Unmarshal(answer, resp)
}

// returns string
// todo refactor dude
func MakeReqToDialogFlow(logger *zap.Logger, req usersdata.DialogFlowReq) (resp usersdata.DialogFlowResp) {
	reqdata, err := json.Marshal(req)
	if err != nil {
		logger.Error("marshal error", zap.Error(err), zap.Any("req", req))
		return
	}

	respdata, err := MakeHttpReq(cfg.DialogFlowURL+"/getAnswer", "POST", reqdata)
	if err != nil {
		logger.Error("make http error dialogflow /getAnswer", zap.Error(err), zap.Any("req", req))
		return
	}

	if err := json.Unmarshal(respdata, &resp); err != nil {
		logger.Error("unmarshal error", zap.Error(err), zap.Any("data", string(respdata)))
		return resp
	}

	if resp.Err != "" {
		logger.Error("got error from dialogflow", zap.Any("response", resp), zap.Any("request", resp))
		return resp
	}
	return resp
}

// MakeReqToTost - makes req to tost service
func MakeReqToTost(method string, req, resp interface{}) (err error) {
	path := cfg.TostURL + method
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	var answer []byte
	switch method {
	case tostcfg.GetRandomTostURL:
		answer, err = MakeHttpReq(path, "GET", data)
	default:
		err = fmt.Errorf("no such method")
	}
	if err != nil {
		return
	}
	return json.Unmarshal(answer, resp)
}

// MakeHttpReq - func for sending http req with given path, method(get or post!) and data
func MakeHttpReq(path, method string, data []byte) (answer []byte, err error) {
	var resp *http.Response
	switch method {
	case "GET":
		resp, err = http.Get(path)
	case "POST":
		resp, err = http.Post(path, "application/json", bytes.NewReader(data))
	default:
		err = fmt.Errorf("method not supported, use get or post methods")
	}
	if err != nil {
		return nil, err
	}

	answer, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}
