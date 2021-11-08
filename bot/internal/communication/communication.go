package communication

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	den4ikdata "github.com/supperdoggy/superSecretDevelopement/structs/request/den4ik"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	den4ikcfg "github.com/supperdoggy/superSecretDevelopement/structs/services/den4ik"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

// MakeUserHttpReq - method handler for users req
// resp must be a pointer!!!!
func MakeUserHttpReq(method string, req, resp interface{}) (err error) {
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	path := cfg.UserURL + method
	var answer []byte
	switch method {
	case cfg.AddFlowerURL:
		answer, err = MakeHttpReq(path, "POST", data)
	case cfg.DialogFlowHandlerURL:
		answer, err = MakeHttpReq(path, "POST", data)
	case cfg.MyFlowersURL:
		answer, err = MakeHttpReq(path, "POST", data)
	case cfg.GiveFlowerURL:
		answer, err = MakeHttpReq(path, "POST", data)
	case cfg.FlowertopURL:
		answer, err = MakeHttpReq(path, "POST", data)
	case cfg.GetFortuneURL:
		answer, err = MakeHttpReq(path, "POST", data)
	case cfg.GetRandomNHIEURL:
		answer, err = MakeHttpReq(path, "POST", data)
	case cfg.GetRandomTostURL:
		answer, err = MakeHttpReq(path, "POST", data)
	case cfg.GetRandomAnekURL:
		answer, err = MakeHttpReq(path, "POST", data)
	default:
		err = fmt.Errorf("no such method")
	}
	if err != nil {
		return
	}
	return json.Unmarshal(answer, resp)
}

func UpdateUser(logger *zap.Logger, usermsg, botmsg *telebot.Message) {
	if usermsg == nil || botmsg == nil {
		return
	}
	var req structs.User = structs.User{
		Telebot: *usermsg.Sender,
		Chats: []structs.Chat{{
			Telebot:    *usermsg.Chat,
			LastOnline: time.Now().Unix(),
		}},
		MessagesUserSent:   []telebot.Message{*usermsg},
		MessagesZhannaSent: []telebot.Message{*botmsg},
	}
	var resp usersdata.AddOrUpdateUserResp
	data, err := json.Marshal(req)
	if err != nil {
		logger.Error("error unmarshalling", zap.Error(err))
		return
	}
	respdata, err := MakeHttpReq(cfg.UserURL+cfg.AddOrUpdateUserURL, "POST", data)
	if err != nil {
		logger.Error("error making req", zap.Error(err))
		return
	}
	err = json.Unmarshal(respdata, &resp)
	if err != nil {
		logger.Error("unmarshal error", zap.Error(err), zap.String("body", string(respdata)))
		return
	}
	if !resp.OK {
		fmt.Printf("UserUpdate error: %+v\n", resp)
	}
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

// grow flower
// TODO REFACTOR
func MakeFlowerReq(id int, chatId int64, langcode string) (msg string, err error) {
	var req = usersdata.FlowerReq{ID: id, NonDying: chatId == int64(Cfg.EdemID)}
	var resp usersdata.FlowerResp

	marshaled, err := json.Marshal(req)
	if err != nil {
		return "communication error", err
	}
	respdata, err := MakeHttpReq(cfg.UserURL+cfg.FlowerURL, "POST", marshaled)
	if err != nil {
		return "communication error", err
	}

	if err := json.Unmarshal(respdata, &resp); err != nil {
		return "communication error", err
	}
	// making request to my flowers to get total and last
	myflowersReq := usersdata.MyFlowersReq{ID: id}
	var myflowersResp usersdata.MyFlowersResp
	// getting total and last
	err = MakeUserHttpReq(cfg.MyFlowersURL, myflowersReq, &myflowersResp)
	if err != nil {
		return "communication error", err
	}

	var replymsg string
	if myflowersResp.Err != "" {
		return "communiction error", errors.New(myflowersResp.Err)
	} else {
		replymsg = fmt.Sprintf(localization.GetLoc("flower_already_have", langcode), myflowersResp.Total, myflowersResp.Last)
	}

	if resp.Err == "cant grow flower" {
		return localization.GetLoc("already_grew_flowers", langcode) + replymsg, nil
	}

	if resp.Err != "" {
		return "communication error", err
	}
	if resp.Dead {
		return fmt.Sprintf(localization.GetLoc("flower_died", langcode)) + replymsg, nil
	}
	if resp.HP == 100 {
		return fmt.Sprintf(localization.GetLoc("flower_grew", langcode), resp.Icon) + replymsg, err
	}
	if resp.Grew {
		return fmt.Sprintf(localization.GetLoc("flower_grew_not_fully", langcode)+replymsg, resp.Up, resp.Extra), err
	}
	return "its not time, try again later...", err
}

func MakeAdminHTTPReq(method string, req, resp interface{}) (err error) {
	path := Cfg.UsersAdminURL + method
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	var dataresp []byte
	switch method {
	case cfg.IsAdminURL:
		dataresp, err = MakeHttpReq(path, "POST", data)
	case cfg.ChangeAdminURL:
		dataresp, err = MakeHttpReq(path, "POST", data)
	case cfg.GetAllFlowerTypesURL:
		dataresp, err = MakeHttpReq(path, "GET", nil)
	case cfg.RemoveFlowerURL:
		dataresp, err = MakeHttpReq(path, "POST", data)
	case cfg.AddUserFlowerURL:
		dataresp, err = MakeHttpReq(path, "POST", data)
	default:
		err = fmt.Errorf("no such method")
	}
	if err != nil {
		return
	}
	return json.Unmarshal(dataresp, resp)
}

func GetCard(id int) (den4ikdata.GetCardResp, error) {
	var req den4ikdata.GetCardReq
	var resp den4ikdata.GetCardResp
	req.SessionID = id

	data, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}

	dataresp, err := MakeHttpReq(den4ikcfg.URL+den4ikcfg.GetCardURL, "POST", data)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(dataresp, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Err != "" {
		return resp, errors.New(resp.Err)
	}
	return resp, nil
}

func ResetDen4ik(sessionID int) (resp den4ikdata.ResetSessionResp, err error) {
	var req = den4ikdata.ResetSessionReq{SessionID: sessionID}

	data, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}

	dataresp, err := MakeHttpReq(den4ikcfg.URL+den4ikcfg.SessionReset, "POST", data)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(dataresp, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Err != "" {
		return resp, errors.New(resp.Err)
	}
	return resp, nil
}
