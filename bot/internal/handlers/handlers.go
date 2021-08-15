package handlers

import (
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
	err := communication.MakeUserHttpReq(cfg.GetFortuneURL, req, &resp)
	if err != nil {
		log.Println("error making request")
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
	var req = usersdata.GetRandomAnekReq{ID: m.Sender.ID}
	var resp usersdata.GetRandomAnekResp
	err := communication.MakeUserHttpReq(cfg.GetRandomAnekURL, req, &resp)
	if err != nil {
		log.Println("handlers.go -> anek() -> make req error:", err.Error())
		return
	}
	botmsg, err := h.Bot.Reply(m, resp.Text)
	if err != nil {
		log.Println("handlers.go -> anek() -> reply error:", err.Error())
		return
	}
	go communication.UpdateUser(m, botmsg)
}

func (h *Handlers) Tost(m *telebot.Message) {
	var req = usersdata.GetRandomTostReq{ID: m.Sender.ID}
	var resp usersdata.GetRandomTostResp
	err := communication.MakeUserHttpReq(cfg.GetRandomTostURL, req, &resp)
	if err != nil {
		log.Println("handlers.go -> tost() -> make req error:", err.Error())
		return
	}
	botmsg, err := h.Bot.Reply(m, resp.Text)
	if err != nil {
		log.Println("handlers.go -> tost() -> reply error:", err.Error())
		return
	}
	go communication.UpdateUser(m, botmsg)
}

// todo onion architecture here
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

	var req = usersdata.DialogFlowReq{ID: types.String(m.Sender.ID), Text: m.Text}
	var resp usersdata.DialogFlowResp
	err := communication.MakeUserHttpReq(cfg.DialogFlowHandlerURL, req, &resp)
	if err != nil {
		log.Println("onTextHandler() -> req error:", err.Error())
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
	err := communication.MakeUserHttpReq(cfg.MyFlowersURL, req, &resp)
	if err != nil {
		log.Println("myflowers() -> MakeUserHttpReq(myflowers) err:", err.Error())
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
	receiver := m.ReplyTo.Sender

	var req = usersdata.GiveFlowerReq{
		Owner:    m.Sender.ID,
		Reciever: receiver.ID,
		Last:     true,
	}
	var resp usersdata.GiveFlowerResp
	err := communication.MakeUserHttpReq(cfg.GiveFlowerURL, req, &resp)
	if err != nil {
		log.Printf("handlers.go -> user give req error: %v, data:%v\n", err.Error(), req)
		return
	}

	if resp.Err != "" {
		log.Println(resp.Err)
	}
	var user = receiver.FirstName
	if receiver.Username != "" {
		user = receiver.Username
	}
	b, _ := h.Bot.Reply(m, fmt.Sprintf(localization.GetLoc("give_flower_good"), user, resp.Flower.Name+" "+resp.Flower.Icon))
	go communication.UpdateUser(m, b)
}

// Flowertop - forms user top by total amount of flowers
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
	err := communication.MakeUserHttpReq(cfg.FlowertopURL, req, &resp)
	if err != nil {
		log.Printf("handlers.go -> flowertop() -> MakeUserHttpReq('flowertop') error: %v, chatid: %v\n", err.Error(), m.Chat.ID)
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
	err := communication.MakeUserHttpReq(cfg.GetRandomNHIEURL, nil, &resp)
	if err != nil {
		h.Bot.Reply(m, localization.GetLoc("error"))
		return
	}

	h.Bot.Reply(m, resp.Result.Text)
}
