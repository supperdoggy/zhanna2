package handlers

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	service "github.com/supperdoggy/superSecretDevelopement/bot/internal/service"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/night-codes/types.v1"
	"gopkg.in/tucnak/telebot.v2"
)

type Handlers struct {
	bot     *telebot.Bot
	service service.Service
	logger  *zap.Logger
}

func NewHandlers(b *telebot.Bot, s service.Service, l *zap.Logger) *Handlers {
	return &Handlers{
		bot:     b,
		service: s,
		logger:  l,
	}
}

// Start - handles /start command and sends text response
func (h *Handlers) Start(m *telebot.Message) {
	var response string
	// todo: create id checker and answer variations for different users
	response = localization.GetLoc("prod_welcome")
	botmsg, err := h.bot.Reply(m, response)
	if err != nil {
		h.logger.Error("error replying to message", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

func (h *Handlers) FortuneCookie(m *telebot.Message) {
	var resp usersdata.GetFortuneResp
	req := usersdata.GetFortuneReq{ID: m.Sender.ID}
	err := communication.MakeUserHttpReq(cfg.GetFortuneURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	msg := resp.Fortune.Text
	if resp.Err != "" {
		msg = fmt.Sprintf(localization.GetLoc("fortune"), resp.Err, resp.Fortune.Text)
	}

	botmsg, err := h.bot.Reply(m, msg)
	if err != nil {
		h.logger.Error("error sending answer, FortuneCookie:", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

// anek() - handles /anek command and sends anek text response
func (h *Handlers) Anek(m *telebot.Message) {
	var req = usersdata.GetRandomAnekReq{ID: m.Sender.ID}
	var resp usersdata.GetRandomAnekResp
	err := communication.MakeUserHttpReq(cfg.GetRandomAnekURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	botmsg, err := h.bot.Reply(m, resp.Text)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

func (h *Handlers) Tost(m *telebot.Message) {
	var req = usersdata.GetRandomTostReq{ID: m.Sender.ID}
	var resp usersdata.GetRandomTostResp
	err := communication.MakeUserHttpReq(cfg.GetRandomTostURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	botmsg, err := h.bot.Reply(m, resp.Text)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

// todo onion architecture here
func (h *Handlers) Flower(m *telebot.Message) {
	replymsg, err := communication.MakeFlowerReq(m.Sender.ID, m.Chat.ID)
	if err != nil {
		h.logger.Error("Error making request to flower", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		_, _ = h.bot.Reply(m, localization.GetLoc("error"))
		return
	}

	botmsg, err := h.bot.Reply(m, replymsg)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
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
		h.logger.Error("Error making request to users", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}

	if resp.Err != "" {
		h.logger.Error("got error in resp", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	botmsg, err := h.bot.Reply(m, resp.Answer)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

func (h *Handlers) MyFlowers(m *telebot.Message) {
	var req = usersdata.MyFlowersReq{ID: m.Sender.ID}
	var resp usersdata.MyFlowersResp
	err := communication.MakeUserHttpReq(cfg.MyFlowersURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}

	if resp.Err != "" {
		h.logger.Error("got error in resp", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.bot.Reply(m, resp.Err)
		return
	}

	var answerstr = fmt.Sprintf(localization.GetLoc("my_flower"), resp.Total, resp.Last)
	for _, v := range resp.Flowers {
		answerstr += fmt.Sprintf("%v - %v\n", v.Name, v.Amount)
	}

	botmsg, err := h.bot.Reply(m, answerstr)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

func (h *Handlers) GiveOneFlower(m *telebot.Message) {
	if !m.IsReply() {
		b, _ := h.bot.Reply(m, localization.GetLoc("give_flower_need_reply"))
		communication.UpdateUser(h.logger, m, b)
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
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}

	if resp.Err != "" {
		h.logger.Error("got error from users", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
	}
	var user = receiver.FirstName
	if receiver.Username != "" {
		user = receiver.Username
	}
	b, err := h.bot.Reply(m, fmt.Sprintf(localization.GetLoc("give_flower_good"), user, resp.Flower.Name+" "+resp.Flower.Icon))
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, b)
}

// Flowertop - forms user top by total amount of flowers
// works only in group chats and supergroups
func (h *Handlers) Flowertop(m *telebot.Message) {
	// check for private chat
	if m.Chat.Type == telebot.ChatPrivate {
		botmsg, _ := h.bot.Reply(m, localization.GetLoc("command_only_in_group"))
		communication.UpdateUser(h.logger, m, botmsg)
		return
	}
	var req = usersdata.FlowertopReq{ChatId: types.Int(m.Chat.ID)}
	var resp usersdata.FlowertopResp
	err := communication.MakeUserHttpReq(cfg.FlowertopURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		botmsg, _ := h.bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(h.logger, m, botmsg)
		return
	}
	var msg = fmt.Sprintf(localization.GetLoc("chat_top"), m.Chat.FirstName+""+m.Chat.LastName)
	for k, v := range resp.Result {
		msg += fmt.Sprintf("%v. %v - %v 🌷\n", k+1, v.Username, v.Total)
	}
	botmsg, err := h.bot.Reply(m, msg)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

// handler for danet, returns agree or disagree message to user
func (h *Handlers) Danet(m *telebot.Message) {
	answer := localization.GetRandomDanet()
	botmsg, err := h.bot.Reply(m, answer)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

func (h *Handlers) Neverhaveiever(m *telebot.Message) {
	var resp usersdata.GetRandomNHIEresp
	err := communication.MakeUserHttpReq(cfg.GetRandomNHIEURL, nil, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.bot.Reply(m, localization.GetLoc("error"))
		return
	}

	botmsg, err := h.bot.Reply(m, resp.Result.Text)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}

func (h *Handlers) Den4ikGame(m *telebot.Message) {
	pics, err := h.service.GetCard(int(m.Chat.ID))
	if err != nil && err != service.ErrSessionEnded {
		if _, err := h.bot.Send(m.Chat, localization.GetLoc("error")); err != nil {
			h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		}
		return
		// check if session is ended
	} else if err == service.ErrSessionEnded {
		if _, err := h.bot.Send(m.Chat, localization.GetLoc("den4ik_game_end")); err != nil {
			h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		}
		return
	}
	for _, v := range pics {
		_, err = h.bot.Send(m.Chat, v)
		if err != nil {
			h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		}
	}
}

func (h *Handlers) ResetDen4ik(m *telebot.Message) {
	msg, err := h.service.ResetDen4ik(int(m.Chat.ID))
	if err != nil && err != mgo.ErrNotFound {
		h.logger.Error("reset den4ik error", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		if _, err := h.bot.Send(m.Chat, localization.GetLoc("error")); err != nil {
			h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		}
		return
	}
	botmsg, err := h.bot.Send(m.Chat, msg)
	if err != nil {
		h.logger.Error("Error replying", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(h.logger, m, botmsg)
}
