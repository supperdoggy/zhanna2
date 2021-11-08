package admin_handlers

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/config"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	flowersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/flowers"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"go.uber.org/zap"
	"gopkg.in/night-codes/types.v1"
	"strconv"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

type AdminHandlers struct {
	bot    *telebot.Bot
	logger *zap.Logger
}

func NewAdminHandlers(bot *telebot.Bot, logger *zap.Logger) *AdminHandlers {
	return &AdminHandlers{
		bot:    bot,
		logger: logger,
	}
}

// botReplyAndSave for replying and saving user message
// TODO CREATE PACKAGE TO NOT TO COPY THAT CODE
func (h *AdminHandlers) botReplyAndSave(m *telebot.Message, what interface{}, options ...interface{}) {
	botmsg, err := h.bot.Reply(m, what)
	if err != nil {
		h.logger.Error("error replying to user message",
			zap.Error(err),
			zap.Any("user", m.Sender),
			zap.Any("chat", m.Chat),
			zap.Any("what", what),
		)
	}
	communication.UpdateUser(h.logger, m, botmsg)

	if what != localization.GetLoc("error", m.Sender.LanguageCode) {
		h.logger.Info("handled user request",
			zap.String("status", "200"),
			zap.Any("user", m.Sender),
			zap.Any("message", m.Text),
			zap.Any("bot response", botmsg.Text))
		return
	}

	h.logger.Info("error handling user request",
		zap.String("status", "400"),
		zap.Any("user", m.Sender),
		zap.Any("message", m.Text),
		zap.Any("bot response", botmsg.Text))

	if !config.GetConfig(h.logger).ErrorAdminNotification {
		return
	}

	// if what is error I send error message to me
	m.Chat.ID = Cfg.NeMoksID
	botmsg, err = h.bot.Send(m.Chat, localization.GetLoc("send_error_to_master",
		"ru",
		m.Sender.Username,
		zap.Any("user", m.Sender).Interface,
		zap.Any("chat", m.Chat).Interface,
		zap.Any("options", options).Interface))
	if err != nil {
		h.logger.Error("error replying to admin message",
			zap.Error(err),
			zap.Any("user", m.Sender),
			zap.Any("chat", m.Chat),
			zap.Any("what", what),
		)
	}
}

// botSendAndSave for sending and saving user message
func (h *AdminHandlers) botSendAndSave(msg *telebot.Message, to telebot.Recipient, what interface{}, options ...interface{}) {
	botmsg, err := h.bot.Send(to, what)
	if err != nil {
		h.logger.Error("error replying to message",
			zap.Error(err),
			zap.Any("user", msg.Sender),
			zap.Any("chat", msg.Chat),
			zap.Any("what", what),
		)
	}
	communication.UpdateUser(h.logger, msg, botmsg)
}

func (ah *AdminHandlers) AdminHelp(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin", m.Sender.LanguageCode))
		return
	}

	text := localization.GetLoc("admin_help", m.Sender.LanguageCode)
	ah.botReplyAndSave(m, text)
}

// addFlower - adds new flower type
func (ah *AdminHandlers) AddFlower(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin", m.Sender.LanguageCode))
		return
	}

	text := strings.Split(m.Text[11:], "-")
	if len(text) != 3 {
		ah.botReplyAndSave(m, localization.GetLoc("add_flower", m.Sender.LanguageCode))
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
		ah.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}

	var msg string = localization.GetLoc("done", m.Sender.LanguageCode)
	if !resp.OK {
		msg = fmt.Sprintf("%s - %s", localization.GetLoc("error", m.Sender.LanguageCode), resp.Err)
	}

	ah.botReplyAndSave(m, msg)
}

func (ah *AdminHandlers) Admin(m *telebot.Message) {
	if m.Sender.ID != Cfg.NeMoksID {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin", m.Sender.LanguageCode))
		return
	}
	if !m.IsReply() || m.ReplyTo.Sender.ID == m.Sender.ID {
		ah.botReplyAndSave(m, localization.GetLoc("need_reply_create_admin", m.Sender.LanguageCode))
		return
	}

	req := usersdata.AdminReq{ID: m.ReplyTo.Sender.ID}
	var resp usersdata.AdminResp
	err := communication.MakeAdminHTTPReq(cfg.ChangeAdminURL, req, &resp)
	if err != nil {
		ah.logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}

	if !resp.OK {
		ah.logger.Error("marsal error", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}

	ah.botReplyAndSave(m, fmt.Sprintf(localization.GetLoc("change_admin", m.Sender.LanguageCode), req.ID, resp.Admin))
}

func (ah *AdminHandlers) AllFlowers(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin", m.Sender.LanguageCode))
		return
	}

	var resp usersdata.GetAllFlowerTypesResp
	err := communication.MakeAdminHTTPReq(cfg.GetAllFlowerTypesURL, nil, &resp)
	if err != nil {
		ah.logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}

	if resp.Err != "" {
		ah.logger.Error("got error from resp", zap.Any("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}

	var text string
	for _, v := range resp.Result {
		text += fmt.Sprintf("%v:%v - %v\n", v.ID, v.Icon, v.Name)
	}
	text += fmt.Sprintf("len %v", len(resp.Result))
	ah.botReplyAndSave(m, text)
}

func (ah *AdminHandlers) RemoveFlower(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin", m.Sender.LanguageCode))
		return
	}

	splitted := strings.Split(m.Text, " ")
	if len(splitted) != 2 {
		return
	}
	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		ah.logger.Error("error converting id", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("remove_flower_need_id", m.Sender.LanguageCode))
		return
	}

	req := usersdata.RemoveFlowerReq{ID: types.Uint64(id)}
	var resp usersdata.RemoveFlowerResp
	err = communication.MakeAdminHTTPReq(cfg.RemoveFlowerURL, req, &resp)
	if err != nil {
		ah.logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}

	if !resp.OK {
		ah.logger.Error("got error from resp", zap.Any("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}
	ah.botReplyAndSave(m, "ok")
}

func (ah *AdminHandlers) CheckAdmin(id int) (bool, error) {
	if id == Cfg.NeMoksID {
		return true, nil
	}
	req := usersdata.IsAdminReq{ID: id}
	var resp usersdata.IsAdminResp
	err := communication.MakeAdminHTTPReq(cfg.IsAdminURL, req, &resp)
	if err != nil {
		ah.logger.Error("Error making admin request to user", zap.Error(err), zap.Int("user_id", id))
		return false, err
	}

	if resp.Err != "" {
		ah.logger.Error("got error from resp", zap.Any("error", resp.Err), zap.Int("user_id", id))
		return false, fmt.Errorf(resp.Err)
	}
	return resp.Result, nil
}

func (ah *AdminHandlers) AddUserFlowerRandom(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin", m.Sender.LanguageCode))
		return
	}

	if !m.IsReply() {
		ah.botReplyAndSave(m, localization.GetLoc("need_reply_add_flower", m.Sender.LanguageCode))
		return
	}

	var err error
	var req flowersdata.AddUserFlowerReq
	req.UserID = m.Sender.ID
	req.RandomFlower = true
	splitted := strings.Split(m.Text, " ")
	if len(splitted) != 1 {
		req.Count, err = strconv.Atoi(splitted[1])
		if err != nil {
			ah.botReplyAndSave(m, localization.GetLoc("add_flower_args", m.Sender.LanguageCode))
			return
		}
		if req.Count < 2 {
			ah.botReplyAndSave(m, localization.GetLoc("count_error", m.Sender.LanguageCode))
			return
		}
		req.Multiple = true
	}

	var resp flowersdata.AddUserFlowerResp
	err = communication.MakeAdminHTTPReq(cfg.AddUserFlowerURL, req, &resp)
	if err != nil {
		ah.logger.Error("error making request to admin", zap.Error(err), zap.Any("req", req))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}

	if resp.Error != "" {
		ah.logger.Error("got error from flowers", zap.Any("resp", resp), zap.Any("req", req))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}
	name := m.ReplyTo.Sender.Username
	if name == "" {
		name = m.ReplyTo.Sender.FirstName + " " + m.ReplyTo.Sender.LastName
	}
	if req.Multiple {
		ah.botReplyAndSave(m, localization.GetLoc("add_user_flower_multiple", m.Sender.LanguageCode, name, len(resp.Flowers)))
		return
	}
	ah.botReplyAndSave(m, localization.GetLoc("add_user_flower", m.Sender.LanguageCode, resp.Flowers[0].Icon, name))

}

func (ah *AdminHandlers) AddUserFlowerByID(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin", m.Sender.LanguageCode))
		return
	}

	if !m.IsReply() {
		ah.botReplyAndSave(m, localization.GetLoc("need_reply_add_flower", m.Sender.LanguageCode))
		return
	}

	var err error
	var req flowersdata.AddUserFlowerReq
	req.UserID = m.ReplyTo.Sender.ID
	splitted := strings.Split(m.Text, " ")
	if len(splitted) != 1 {
		flowerid, err := strconv.Atoi(splitted[1])
		if err != nil {
			ah.botReplyAndSave(m, localization.GetLoc("add_flower_args", m.Sender.LanguageCode))
			return
		}
		req.FlowerID = uint64(flowerid)
	}

	var resp flowersdata.AddUserFlowerResp
	err = communication.MakeAdminHTTPReq(cfg.AddUserFlowerURL, req, &resp)
	if err != nil {
		ah.logger.Error("error making request to admin", zap.Error(err), zap.Any("req", req))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}

	if resp.Error != "" {
		ah.logger.Error("got error from flowers", zap.Any("resp", resp), zap.Any("req", req))
		ah.botReplyAndSave(m, localization.GetLoc("error", m.Sender.LanguageCode))
		return
	}
	name := m.ReplyTo.Sender.Username
	if name == "" {
		name = m.ReplyTo.Sender.FirstName + " " + m.ReplyTo.Sender.LastName
	}
	ah.botReplyAndSave(m, localization.GetLoc("add_user_flower", m.Sender.LanguageCode, resp.Flowers[0].Icon, name))

}
