package admin_handlers

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
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
func (h *AdminHandlers) botReplyAndSave(m *telebot.Message, what interface{}, options ...interface{}) {
	botmsg, err := h.bot.Reply(m, what)
	if err != nil {
		h.logger.Error("error replying to message",
			zap.Error(err),
			zap.Any("user", m.Sender),
			zap.Any("chat", m.Chat),
			zap.Any("what", what),
		)
	}
	communication.UpdateUser(h.logger, m, botmsg)
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
		ah.botReplyAndSave(m, localization.GetLoc("not_admin"))
		return
	}

	text := localization.GetLoc("admin_help")
	ah.botReplyAndSave(m, text)
}

// addFlower - adds new flower type
func (ah *AdminHandlers) AddFlower(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin"))
		return
	}

	text := strings.Split(m.Text[11:], "-")
	if len(text) != 3 {
		ah.botReplyAndSave(m, localization.GetLoc("add_flower"))
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
		ah.botReplyAndSave(m, localization.GetLoc("error"))
		return
	}

	var msg string = localization.GetLoc("done")
	if !resp.OK {
		msg = fmt.Sprintf("%s - %s", localization.GetLoc("error"), resp.Err)
	}

	ah.botReplyAndSave(m, msg)
}

func (ah *AdminHandlers) Admin(m *telebot.Message) {
	if m.Sender.ID != Cfg.NeMoksID {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin"))
		return
	}
	if !m.IsReply() || m.ReplyTo.Sender.ID == m.Sender.ID {
		ah.botReplyAndSave(m, localization.GetLoc("need_reply"))
		return
	}

	req := usersdata.AdminReq{ID: m.ReplyTo.Sender.ID}
	var resp usersdata.AdminResp
	err := communication.MakeAdminHTTPReq(cfg.ChangeAdminURL, req, &resp)
	if err != nil {
		ah.logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error"))
		return
	}

	if !resp.OK {
		ah.logger.Error("marsal error", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error"))
		return
	}

	ah.botReplyAndSave(m, fmt.Sprintf(localization.GetLoc("change_admin"), req.ID, resp.Admin))
}

func (ah *AdminHandlers) AllFlowers(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		ah.botReplyAndSave(m, localization.GetLoc("not_admin"))
		return
	}

	var resp usersdata.GetAllFlowerTypesResp
	err := communication.MakeAdminHTTPReq(cfg.GetAllFlowerTypesURL, nil, &resp)
	if err != nil {
		ah.logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error"))
		return
	}

	if resp.Err != "" {
		ah.logger.Error("got error from resp", zap.Any("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error"))
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
		ah.botReplyAndSave(m, localization.GetLoc("not_admin"))
		return
	}

	splitted := strings.Split(m.Text, " ")
	if len(splitted) != 2 {
		return
	}
	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		ah.logger.Error("error converting id", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("remove_flower_need_id"))
		return
	}

	req := usersdata.RemoveFlowerReq{ID: types.Uint64(id)}
	var resp usersdata.RemoveFlowerResp
	err = communication.MakeAdminHTTPReq(cfg.RemoveFlowerURL, req, &resp)
	if err != nil {
		ah.logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error"))
		return
	}

	if !resp.OK {
		ah.logger.Error("got error from resp", zap.Any("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.botReplyAndSave(m, localization.GetLoc("error"))
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
