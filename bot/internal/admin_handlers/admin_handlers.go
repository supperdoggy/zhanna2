package admin_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
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

	text := "/admin - set/unset admin\n" +
		"/addFlower - add new flower type\n" +
		"/removeFlower - remove flower type\n" +
		"/allFlowers - returns flower types list\n"
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
	data := obj{"icon": text[0], "name": text[1], "type": text[2]}
	_, err := communication.MakeUserHttpReq(cfg.AddFlowerURL, data)
	if err != nil {
		log.Println("admin_handlers.go -> addFlower() -> MakeUserHttpReq error:", err.Error())
		return
	}
	botmsg, _ := ah.Bot.Reply(m, "Done!")
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

	data, err := communication.MakeAdminHTTPReq(cfg.ChangeAdminURL, obj{"id": m.ReplyTo.Sender.ID})
	if err != nil {
		log.Printf("admin_handlers.go -> admin() -> MakeAdminHTTPReq error: %v id: %v\n", err.Error(), m.ReplyTo.Sender.ID)
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(m, botmsg)
		return
	}

	var resp struct {
		Err   string `json:"err"`
		Admin bool   `json:"admin"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil || resp.Err != "" {
		log.Printf("admin_handlers.go -> admin() -> Marshal error: %v body: %v, resp error: %v\n", err.Error(), string(data), resp.Err)
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(m, botmsg)
		return
	}

	// todo localization
	botmsg, _ := ah.Bot.Reply(m, fmt.Sprintf("Пользователь %v admin: %v\n", m.ReplyTo.Sender.ID, resp.Admin))
	go communication.UpdateUser(m, botmsg)

	bot.Send(m.Sender, fmt.Sprintf("Пользователь %v admin: %v\n", m.ReplyTo.Sender.ID, resp.Admin))
	return
}

func (ah *AdminHandlers) AllFlowers(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(m, botmsg)
		return
	}

	data, err := communication.MakeAdminHTTPReq(cfg.GetAllFlowerTypesURL, obj{})
	if err != nil {
		log.Println("admin_handlers.go -> allFlowers() -> error makin req err:", err.Error())
		return
	}

	var resp struct {
		Result []structs.Flower `json:"result"`
		Err    string           `json:"err"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Println("admin_handlers.go -> allFlowers() -> Unmarshal err:", err.Error(), string(data))
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

	data, err := communication.MakeAdminHTTPReq(cfg.RemoveFlowerURL, obj{"id": id})
	if err != nil {
		log.Println("admin_handlers.go -> removeFlower() -> removeFlower err:", err.Error())
		ah.Bot.Reply(m, localization.GetLoc("error"))
		return
	}
	var resp struct {
		Err string `json:"err"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Println("admin_handlers.go -> removeFlower() -> Unmarshal err:", err.Error(), string(data))
		ah.Bot.Reply(m, localization.GetLoc("error"))
		return
	}
	if resp.Err != "" {
		log.Println("admin_handlers.go -> removeFlower() -> resp err:", resp)
		return
	}
	ah.Bot.Reply(m, "ok")
}

func (ah *AdminHandlers) CheckAdmin(id int) (bool, error) {
	if id == Cfg.NeMoksID {
		return true, nil
	}
	data, err := communication.MakeAdminHTTPReq(cfg.IsAdminURL, obj{"id": id})
	if err != nil {
		log.Println("admin_auth.go -> checkAdmin() -> isAdmin method req error:", err)
		return false, err
	}
	var resp struct {
		Err    string `json:"err"`
		Result bool   `json:"result"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return false, err
	}

	if resp.Err != "" {
		log.Println("admin_auth.go -> checkAdmin() -> resp error:", resp.Err)
		return false, fmt.Errorf(resp.Err)
	}
	return resp.Result, nil
}
