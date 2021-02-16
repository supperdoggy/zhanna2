package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

func adminHelp(m *telebot.Message) {
	if admin, err := checkAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := bot.Reply(m, getLoc("not_admin"))
		UpdateUser(m, botmsg)
		return
	}

	text := "/admin - set/unset admin\n" +
		"/addFlower - add new flower type\n" +
		"/removeFlower - remove flower type\n" +
		"/allFlowers - returns flower types list\n"
	bot.Reply(m, text)
}

// addFlower - adds new flower type
func addFlower(m *telebot.Message) {
	if admin, err := checkAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := bot.Reply(m, getLoc("not_admin"))
		UpdateUser(m, botmsg)
		return
	}

	text := split(m.Text[11:], "-")
	if len(text) != 3 {
		bmsg, _ := bot.Reply(m, getLoc("add_flower"))
		go UpdateUser(m, bmsg)
		return
	}
	data := obj{"icon": text[0], "name": text[1], "type": text[2]}
	_, err := MakeUserHttpReq("addFlower", data)
	if err != nil {
		log.Println("admin_handlers.go -> addFlower() -> MakeUserHttpReq error:", err.Error())
		return
	}
	botmsg, _ := bot.Reply(m, "Done!")
	go UpdateUser(m, botmsg)
}

func admin(m *telebot.Message) {
	if m.Sender.ID != NeMoksID {
		botmsg, _ := bot.Reply(m, getLoc("not_admin"))
		UpdateUser(m, botmsg)
		return
	}
	if !m.IsReply() || m.ReplyTo.Sender.ID == m.Sender.ID {
		botmsg, _ := bot.Reply(m, getLoc("need_reply"))
		UpdateUser(m, botmsg)
		return
	}

	data, err := MakeAdminHTTPReq("admin", obj{"id": m.ReplyTo.Sender.ID})
	if err != nil {
		log.Printf("admin_handlers.go -> admin() -> MakeAdminHTTPReq error: %v id: %v\n", err.Error(), m.ReplyTo.Sender.ID)
		botmsg, _ := bot.Reply(m, getLoc("error"))
		UpdateUser(m, botmsg)
		return
	}

	var resp struct {
		Err   string `json:"err"`
		Admin bool   `json:"admin"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil || resp.Err != "" {
		log.Printf("admin_handlers.go -> admin() -> Marshal error: %v body: %v, resp error: %v\n", err.Error(), string(data), resp.Err)
		botmsg, _ := bot.Reply(m, getLoc("error"))
		UpdateUser(m, botmsg)
		return
	}

	botmsg, _ := bot.Reply(m, fmt.Sprintf("Пользователь %v admin: %v\n", m.ReplyTo.Sender.ID, resp.Admin))
	go UpdateUser(m, botmsg)

	bot.Send(m.Sender, fmt.Sprintf("Пользователь %v admin: %v\n", m.ReplyTo.Sender.ID, resp.Admin))
	return
}

func allFlowers(m *telebot.Message) {
	if admin, err := checkAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := bot.Reply(m, getLoc("not_admin"))
		UpdateUser(m, botmsg)
		return
	}

	data, err := MakeAdminHTTPReq("getAllFlowerTypes", obj{})
	if err != nil {
		log.Println("admin_handlers.go -> allFlowers() -> error makin req err:", err.Error())
		return
	}

	var resp struct {
		Result []Flower `json:"result"`
		Err    string   `json:"err"`
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
	botmsg, _ := bot.Reply(m, text)
	UpdateUser(m, botmsg)
}

func removeFlower(m *telebot.Message) {
	if admin, err := checkAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := bot.Reply(m, getLoc("not_admin"))
		UpdateUser(m, botmsg)
		return
	}

	splitted := strings.Split(m.Text, " ")
	if len(splitted) != 2 {
		return
	}
	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		log.Println("admin_handlers.go -> removeFlower() -> Atoi err:", err.Error())
		bot.Reply(m, "error get id, need /removeFlower <id>")
		return
	}

	data, err := MakeAdminHTTPReq("removeFlower", obj{"id": id})
	if err != nil {
		log.Println("admin_handlers.go -> removeFlower() -> removeFlower err:", err.Error())
		bot.Reply(m, getLoc("error"))
		return
	}
	var resp struct {
		Err string `json:"err"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Println("admin_handlers.go -> removeFlower() -> Unmarshal err:", err.Error(), string(data))
		bot.Reply(m, getLoc("error"))
		return
	}
	if resp.Err != "" {
		log.Println("admin_handlers.go -> removeFlower() -> resp err:", resp)
		return
	}
	bot.Reply(m, "ok")
}
