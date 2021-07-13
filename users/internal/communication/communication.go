package communication

import (
	"bytes"
	"encoding/json"
	"fmt"
	anekscfg "github.com/supperdoggy/superSecretDevelopement/structs/services/aneks"
	flowerscfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	tostcfg "github.com/supperdoggy/superSecretDevelopement/structs/services/tost"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"io/ioutil"
	"net/http"
)

type obj map[string]interface{}

// MakeReqToAnek - makes req to anek service
func MakeReqToAnek(method string, data []byte) (answer []byte, err error) {
	path := cfg.AnekURL + method
	switch method {
	case anekscfg.GetRandomAnekURL:
		answer, err = MakeHttpReq(path, "GET", data)
	default:
		err = fmt.Errorf("no such method")
	}
	return
}

// MakeReqToFlowers - makes req to flowers service
func MakeReqToFlowers(method string, data interface{}) (answer []byte, err error) {
	path := cfg.FlowersURL + method
	reqData, err := json.Marshal(data)
	if err != nil {
		return
	}
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
	default:
		err = fmt.Errorf("no such method")
	}
	return
}

// returns string
func MakeReqToDialogFlow(message string) (answer string, err error) {
	req, err := json.Marshal(obj{"message": message})
	if err != nil {
		fmt.Println("MakeReqToDialogFlow() -> json.Marshal error:", err.Error())
		return
	}

	resp, err := MakeHttpReq(cfg.DialogFlowURL+"/getAnswer", "POST", req)
	if err != nil {
		fmt.Println("MakeReqToDialogFlow() -> makeHttpReq(/getAnswer) error:", err.Error())
		return
	}

	var respStruct struct {
		Answer string `json:"answer"`
		Err    string `json:"err"`
	}

	if err := json.Unmarshal(resp, &respStruct); err != nil {
		fmt.Println("MakeReqToDialogFlow() -> unmarshal error:", err.Error())
		return "", err
	}

	if respStruct.Err != "" {
		fmt.Println("MakeReqToDialogFlow() -> got an error from dialogflow:", respStruct.Err)
		return "", fmt.Errorf(respStruct.Err)
	}
	return respStruct.Answer, nil
}

// MakeReqToTost - makes req to tost service
func MakeReqToTost(method string, data []byte) (answer []byte, err error) {
	path := cfg.TostURL + method
	switch method {
	case tostcfg.GetRandomTostURL:
		answer, err = MakeHttpReq(path, "GET", data)
	default:
		err = fmt.Errorf("no such method")
	}
	return
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
