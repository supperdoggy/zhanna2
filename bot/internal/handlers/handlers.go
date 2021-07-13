package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"log"

	"gopkg.in/tucnak/telebot.v2"
)

type Handlers struct {
	Bot *telebot.Bot
}

// Start - handles /start command and sends text response
// todo below
func (h *Handlers) Start(m *telebot.Message) {
	var response string
	// todo: create id checker and answer variations for different users
	response = localization.GetLoc("prod_welcome")
	botmsg, err := h.Bot.Reply(m, response)
	if err != nil {
		log.Println("handlers.go -> start() -> error:", err.Error(), ", user id:", m.Sender.ID)
		return
	}
	go communication.UpdateUser(m, botmsg)
}

// TODO:
func (h *Handlers) FortuneCookie(m *telebot.Message) {
	var resp struct {
		Fortune structs.Cookie `json:"fortune"`
		Err     string         `json:"err"`
	}
	data := obj{"id": m.Sender.ID}
	r, err := communication.MakeUserHttpReq(cfg.GetFortuneURL, data)
	if err != nil {
		log.Println("error making request")
		return
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		log.Println("error unmarshalling")
		return
	}
	msg := resp.Fortune.Text
	if resp.Err != "" {
		msg = fmt.Sprintf("%v\n\n%v", resp.Err, resp.Fortune.Text)
	}

	botmsg, err := h.Bot.Reply(m, msg)
	if err != nil {
		log.Println("error sending answer, FortuneCookie:", err.Error())
		return
	}
	go communication.UpdateUser(m, botmsg)
}

// anek() - handles /anek command and sends anek text response
func (h *Handlers) Anek(m *telebot.Message) {
	anekAnswer, err := communication.MakeRandomAnekHttpReq(m.Sender.ID)
	if err != nil {
		log.Println("handlers.go -> anek() -> make req error:", err.Error())
		return
	}
	botmsg, err := h.Bot.Reply(m, anekAnswer.Text)
	if err != nil {
		log.Println("handlers.go -> anek() -> reply error:", err.Error())
		return
	}
	go communication.UpdateUser(m, botmsg)
}

func (h *Handlers) Tost(m *telebot.Message) {
	answerTost, err := communication.MakeRandomTostHttpReq(m.Sender.ID)
	if err != nil {
		log.Println("handlers.go -> tost() -> make req error:", err.Error())
		return
	}
	botmsg, err := h.Bot.Reply(m, answerTost.Text)
	if err != nil {
		log.Println("handlers.go -> tost() -> reply error:", err.Error())
		return
	}
	go communication.UpdateUser(m, botmsg)
}

func (h *Handlers) Flower(m *telebot.Message) {
	resp, err := communication.MakeFlowerReq(m.Sender.ID, m.Chat.ID)
	if err != nil {
		log.Println("handlers.go -> flower() -> MakeFlowerReq() error", err.Error(), m.Sender.ID)
		_, _ = h.Bot.Reply(m, "error occured, contact owner")
		return
	}

	req := obj{"id": m.Sender.ID}
	// getting total and last
	data, err := communication.MakeUserHttpReq(cfg.MyFlowersURL, req)
	if err != nil {
		log.Println("handlers.go -> flower() -> myflowers error:", err.Error())
	} else {
		var respstr struct {
			Total int   `json:"total"`
			Last  uint8 `json:"last"`
		}
		err := json.Unmarshal(data, &respstr)
		if err == nil {
			resp += fmt.Sprintf("\nУ тебя уже %v🌷 и %v🌱", respstr.Total, respstr.Last)
		}
	}

	botmsg, err := h.Bot.Reply(m, resp)
	if err != nil {
		log.Println("handlers.go -> flower() -> bot.Reply() error", err.Error())
		return
	}
	go communication.UpdateUser(m, botmsg)
}

// onTextHandler - makes req to python service and gets message from apiai
func (h *Handlers) OnTextHandler(m *telebot.Message) {

	// if chat is not private then user must reply bot to get answer
	if m.Chat.Type != telebot.ChatPrivate {
		if !m.IsReply() || m.IsReply() && !(m.ReplyTo.Sender.ID == Cfg.ProdBotID || m.ReplyTo.Sender.ID == Cfg.TestbotId) {
			return
		}
	}

	answer, err := communication.MakeUserHttpReq(cfg.DialogFlowHandlerURL, obj{"id": m.Sender.ID, "text": m.Text})
	if err != nil {
		log.Println("onTextHandler() -> req error:", err.Error())
		return
	}

	var resp struct {
		Answer string `json:"answer"`
		Err    string `json:"err"`
	}
	if err := json.Unmarshal(answer, &resp); err != nil {
		log.Println("onTextHandler() -> Unmarshal error:", err.Error())
		return
	}
	if resp.Err != "" {
		log.Println("onTextHandler() -> got error in response:", resp.Err)
		return
	}
	botmsg, _ := h.Bot.Reply(m, resp.Answer)
	go communication.UpdateUser(m, botmsg)
}

func (h *Handlers) MyFlowers(m *telebot.Message) {
	answer, err := communication.MakeUserHttpReq(cfg.MyFlowersURL, obj{"id": m.Sender.ID})
	if err != nil {
		log.Println("myflowers() -> MakeUserHttpReq(myflowers) err:", err.Error())
		return
	}
	var resp struct {
		Flowers map[string]int `json:"flowers"`
		Last    uint8          `json:"last"`
		Total   int            `json:"total"`
		Err     string         `json:"err"`
	}

	if err := json.Unmarshal(answer, &resp); err != nil {
		log.Println("myflowers() -> unmarshal error:", err.Error(), string(answer))
		return
	}

	if resp.Err != "" {
		log.Println("myflowers() -> got error resp from service:", resp.Err)
		h.Bot.Reply(m, resp.Err)
		return
	}

	var answerstr = fmt.Sprintf(localization.GetLoc("my_flower"), resp.Total, resp.Last)
	for k, v := range resp.Flowers {
		answerstr += fmt.Sprintf("%v - %v\n", k, v)
	}
	botmsg, _ := h.Bot.Reply(m, answerstr)
	go communication.UpdateUser(m, botmsg)
}

func (h *Handlers) GiveOneFlower(m *telebot.Message) {
	if !m.IsReply() {
		b, _ := h.Bot.Reply(m, localization.GetLoc("give_flower_need_reply"))
		communication.UpdateUser(m, b)
		return
	}

	data := obj{"last": true, "owner": m.Sender.ID, "reciever": m.ReplyTo.Sender.ID}
	answer, err := communication.MakeUserHttpReq(cfg.GiveFlowerURL, data)
	if err != nil {
		log.Printf("handlers.go -> user give req error: %v, data:%v\n", err.Error(), data)
		return
	}
	var resp struct {
		Err string `json:"err"`
	}
	if err := json.Unmarshal(answer, &resp); err != nil {
		log.Printf("handlers.go -> unmarshal error: %v, body: %v\n", err.Error(), string(answer))
		return
	}
	if resp.Err != "" {
		log.Println(resp.Err)
	}
	b, _ := h.Bot.Reply(m, localization.GetLoc("give_flower_good"))
	go communication.UpdateUser(m, b)
}

// forms user top by total amount of flowers
// works only in group chats and supergroups
func (h *Handlers) Flowertop(m *telebot.Message) {
	// check for private chat
	if m.Chat.Type == telebot.ChatPrivate {
		botmsg, _ := h.Bot.Reply(m, localization.GetLoc("command_only_in_group"))
		communication.UpdateUser(m, botmsg)
		return
	}
	answer, err := communication.MakeUserHttpReq(cfg.FlowertopURL, obj{"chatid": m.Chat.ID})
	if err != nil {
		log.Printf("handlers.go -> flowertop() -> MakeUserHttpReq('flowertop') error: %v, chatid: %v\n", err.Error(), m.Chat.ID)
		botmsg, _ := h.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(m, botmsg)
		return
	}
	var resp struct {
		Top []struct {
			Username string `json:"username"`
			Total    int    `json:"total"`
		} `json:"result"`
	}
	err = json.Unmarshal(answer, &resp)
	if err != nil || len(resp.Top) == 0 {
		log.Printf("handlers.go -> flowertop() -> Unmarshal error:%v, body: %v\n", err, string(answer))
		botmsg, _ := h.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(m, botmsg)
		return
	}
	var msg = fmt.Sprintf(localization.GetLoc("chat_top"), m.Chat.FirstName+""+m.Chat.LastName)
	for k, v := range resp.Top {
		msg += fmt.Sprintf("%v. %v - %v 🌷\n", k+1, v.Username, v.Total)
	}
	botmsg, _ := h.Bot.Reply(m, msg)
	communication.UpdateUser(m, botmsg)
}

// handler for danet, returns agree or disagree message to user
func (h *Handlers) Danet(m *telebot.Message) {
	answer := localization.GetRandomDanet()
	botmsg, err := h.Bot.Reply(m, answer)
	if err != nil {
		log.Printf("handlers.go -> danet() -> Reply() error: %v, id: %v\n", err.Error(), m.Sender.ID)
	}
	communication.UpdateUser(m, botmsg)
}

func (h *Handlers) Neverhaveiever(m *telebot.Message) {
	data, err := communication.MakeUserHttpReq(cfg.GetRandomNHIEURL, nil)
	if err != nil {
		h.Bot.Reply(m, localization.GetLoc("error"))
		return
	}
	var resp struct {
		Err    string       `json:"err"`
		Result structs.NHIE `json:"result"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Printf("handlers.go -> neverhaveiever() -> Unmarshal() error: %v, body: %v\n", err.Error(), string(data))
		return
	}
	h.Bot.Reply(m, resp.Result.Text)
}
