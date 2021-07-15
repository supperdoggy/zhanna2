package admin_handlers

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"gopkg.in/night-codes/types.v1"
	"log"
	"strconv"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

type AdminHandlers struct {
	Bot *telebot.Bot
}

func (ah *AdminHandlers) AdminHelp(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(m, botmsg)
		return
	}

	text := localization.GetLoc("admin_help")
	ah.Bot.Reply(m, text)
}

// addFlower - adds new flower type
func (ah *AdminHandlers) AddFlower(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(m, botmsg)
		return
	}

	text := strings.Split(m.Text[11:], "-")
	if len(text) != 3 {
		bmsg, _ := ah.Bot.Reply(m, localization.GetLoc("add_flower"))
		go communication.UpdateUser(m, bmsg)
		return
	}

	req := usersdata.AddFlowerReq{
		Icon: text[0],
		Name: text[1],
		Type: text[2],
	}
	var resp usersdata.AddFlowerResp
	err := communication.MakeUserHttpReq(cfg.AddFlowerURL, req, &resp)
	if err != nil {
		log.Println("admin_handlers.go -> addFlower() -> MakeUserHttpReq error:", err.Error())
		return
	}

	var msg string = localization.GetLoc("done")
	if !resp.OK {
		msg = fmt.Sprintf("%s - %s", localization.GetLoc("error"), resp.Err)
	}

	botmsg, _ := ah.Bot.Reply(m, msg)
	go communication.UpdateUser(m, botmsg)
}

func (ah *AdminHandlers) Admin(m *telebot.Message) {
	if m.Sender.ID != Cfg.NeMoksID {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(m, botmsg)
		return
	}
	if !m.IsReply() || m.ReplyTo.Sender.ID == m.Sender.ID {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("need_reply"))
		communication.UpdateUser(m, botmsg)
		return
	}

	req := usersdata.AdminReq{ID: m.ReplyTo.Sender.ID}
	var resp usersdata.AdminResp
	err := communication.MakeAdminHTTPReq(cfg.ChangeAdminURL, req, &resp)
	if err != nil {
		log.Printf("admin_handlers.go -> admin() -> MakeAdminHTTPReq error: %v id: %v\n", err.Error(), req.ID)
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(m, botmsg)
		return
	}

	if !resp.OK {
		log.Printf("admin_handlers.go -> admin() -> Marshal error: %v body: %v, resp error: %v\n", err, resp.Err)
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(m, botmsg)
		return
	}

	botmsg, _ := ah.Bot.Reply(m, fmt.Sprintf(localization.GetLoc("change_admin"), req.ID, resp.Admin))
	go communication.UpdateUser(m, botmsg)

	ah.Bot.Send(m.Sender, fmt.Sprintf(localization.GetLoc("change_admin"), req.ID, resp.Admin))
	return
}

func (ah *AdminHandlers) AllFlowers(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(m, botmsg)
		return
	}

	var resp usersdata.GetAllFlowerTypesResp
	err := communication.MakeAdminHTTPReq(cfg.GetAllFlowerTypesURL, nil, &resp)
	if err != nil {
		log.Println("admin_handlers.go -> allFlowers() -> error makin req err:", err.Error())
		return
	}

	if resp.Err != "" {
		log.Println("admin_handlers.go -> allFlowers() -> resp error", resp.Err)
		return
	}

	var text string
	for _, v := range resp.Result {
		text += fmt.Sprintf("%v:%v - %v\n", v.ID, v.Icon, v.Name)
	}
	text += fmt.Sprintf("len %v", len(resp.Result))
	botmsg, _ := ah.Bot.Reply(m, text)
	communication.UpdateUser(m, botmsg)
}

func (ah *AdminHandlers) RemoveFlower(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(m, botmsg)
		return
	}

	splitted := strings.Split(m.Text, " ")
	if len(splitted) != 2 {
		return
	}
	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		log.Println("admin_handlers.go -> removeFlower() -> Atoi err:", err.Error())
		ah.Bot.Reply(m, "error get id, need /removeFlower <id>")
		return
	}

	req := usersdata.RemoveFlowerReq{ID: types.Uint64(id)}
	var resp usersdata.RemoveFlowerResp
	err = communication.MakeAdminHTTPReq(cfg.RemoveFlowerURL, req, &resp)
	if err != nil {
		log.Println("admin_handlers.go -> removeFlower() -> removeFlower err:", err.Error())
		ah.Bot.Reply(m, localization.GetLoc("error"))
		return
	}

	if !resp.OK {
		log.Println("admin_handlers.go -> removeFlower() -> resp err:", resp.Err)
		ah.Bot.Reply(m, localization.GetLoc("error"))
		return
	}
	ah.Bot.Reply(m, "ok")
}

func (ah *AdminHandlers) CheckAdmin(id int) (bool, error) {
	if id == Cfg.NeMoksID {
		return true, nil
	}
	req := usersdata.IsAdminReq{ID: id}
	var resp usersdata.IsAdminResp
	err := communication.MakeAdminHTTPReq(cfg.IsAdminURL, req, &resp)
	if err != nil {
		log.Println("admin_auth.go -> checkAdmin() -> isAdmin method req error:", err)
		return false, err
	}

	if resp.Err != "" {
		log.Println("admin_auth.go -> checkAdmin() -> resp error:", resp.Err)
		return false, fmt.Errorf(resp.Err)
	}
	return resp.Result, nil
}
