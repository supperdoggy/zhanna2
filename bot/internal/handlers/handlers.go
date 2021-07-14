package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"gopkg.in/night-codes/types.v1"
	"log"

	"gopkg.in/tucnak/telebot.v2"
)

type Handlers struct {
	Bot *telebot.Bot
}

// Start - handles /start command and sends text response
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

func (h *Handlers) FortuneCookie(m *telebot.Message) {
	var resp usersdata.GetFortuneResp
	req := usersdata.GetFortuneReq{ID: m.Sender.ID}
	r, err := communication.MakeUserHttpReq(cfg.GetFortuneURL, req)
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
		msg = fmt.Sprintf(localization.GetLoc("fortune"), resp.Err, resp.Fortune.Text)
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
	replymsg, err := communication.MakeFlowerReq(m.Sender.ID, m.Chat.ID)
	if err != nil {
		log.Println("handlers.go -> flower() -> MakeFlowerReq() error", err.Error(), m.Sender.ID)
		_, _ = h.Bot.Reply(m, localization.GetLoc("error"))
		return
	}

	botmsg, err := h.Bot.Reply(m, replymsg)
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

	var req = usersdata.DialogFlowReq{ID: m.Sender.ID, Text: m.Text}
	var resp usersdata.DialogFlowResp
	answer, err := communication.MakeUserHttpReq(cfg.DialogFlowHandlerURL, req)
	if err != nil {
		log.Println("onTextHandler() -> req error:", err.Error())
		return
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
	var req = usersdata.MyFlowersReq{ID: m.Sender.ID}
	var resp usersdata.MyFlowersResp
	answer, err := communication.MakeUserHttpReq(cfg.MyFlowersURL, req)
	if err != nil {
		log.Println("myflowers() -> MakeUserHttpReq(myflowers) err:", err.Error())
		return
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

	var req = usersdata.GiveFlowerReq{
		Owner:    m.Sender.ID,
		Reciever: m.ReplyTo.Sender.ID,
		Last:     true,
	}
	var resp usersdata.GiveFlowerResp
	answer, err := communication.MakeUserHttpReq(cfg.GiveFlowerURL, req)
	if err != nil {
		log.Printf("handlers.go -> user give req error: %v, data:%v\n", err.Error(), req)
		return
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
	var req = usersdata.FlowertopReq{ChatId: types.Int(m.Chat.ID)}
	var resp usersdata.FlowertopResp
	answer, err := communication.MakeUserHttpReq(cfg.FlowertopURL, req)
	if err != nil {
		log.Printf("handlers.go -> flowertop() -> MakeUserHttpReq('flowertop') error: %v, chatid: %v\n", err.Error(), m.Chat.ID)
		botmsg, _ := h.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(m, botmsg)
		return
	}

	err = json.Unmarshal(answer, &resp)
	if err != nil || len(resp.Result) == 0 {
		log.Printf("handlers.go -> flowertop() -> Unmarshal error:%v, body: %v\n", err, string(answer))
		botmsg, _ := h.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(m, botmsg)
		return
	}
	var msg = fmt.Sprintf(localization.GetLoc("chat_top"), m.Chat.FirstName+""+m.Chat.LastName)
	for k, v := range resp.Result {
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
	var resp usersdata.GetRandomNHIEresp
	data, err := communication.MakeUserHttpReq(cfg.GetRandomNHIEURL, nil)
	if err != nil {
		h.Bot.Reply(m, localization.GetLoc("error"))
		return
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		log.Printf("handlers.go -> neverhaveiever() -> Unmarshal() error: %v, body: %v\n", err.Error(), string(data))
		return
	}
	h.Bot.Reply(m, resp.Result.Text)
}
