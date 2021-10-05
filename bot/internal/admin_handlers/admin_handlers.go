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
	Bot    *telebot.Bot
	Logger *zap.Logger
}

func (ah *AdminHandlers) AdminHelp(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(ah.Logger, m, botmsg)
		return
	}

	text := localization.GetLoc("admin_help")
	botmsg, err := ah.Bot.Reply(m, text)
	if err != nil {
		ah.Logger.Error("error replying to message", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(ah.Logger, m, botmsg)
}

// addFlower - adds new flower type
func (ah *AdminHandlers) AddFlower(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(ah.Logger, m, botmsg)
		return
	}

	text := strings.Split(m.Text[11:], "-")
	if len(text) != 3 {
		bmsg, _ := ah.Bot.Reply(m, localization.GetLoc("add_flower"))
		go communication.UpdateUser(ah.Logger, m, bmsg)
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
		ah.Logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}

	var msg string = localization.GetLoc("done")
	if !resp.OK {
		msg = fmt.Sprintf("%s - %s", localization.GetLoc("error"), resp.Err)
	}

	botmsg, err := ah.Bot.Reply(m, msg)
	if err != nil {
		ah.Logger.Error("error replying to message", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(ah.Logger, m, botmsg)
}

func (ah *AdminHandlers) Admin(m *telebot.Message) {
	if m.Sender.ID != Cfg.NeMoksID {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(ah.Logger, m, botmsg)
		return
	}
	if !m.IsReply() || m.ReplyTo.Sender.ID == m.Sender.ID {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("need_reply"))
		communication.UpdateUser(ah.Logger, m, botmsg)
		return
	}

	req := usersdata.AdminReq{ID: m.ReplyTo.Sender.ID}
	var resp usersdata.AdminResp
	err := communication.MakeAdminHTTPReq(cfg.ChangeAdminURL, req, &resp)
	if err != nil {
		ah.Logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(ah.Logger, m, botmsg)
		return
	}

	if !resp.OK {
		ah.Logger.Error("marsal error", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("error"))
		communication.UpdateUser(ah.Logger, m, botmsg)
		return
	}

	botmsg, err := ah.Bot.Reply(m, fmt.Sprintf(localization.GetLoc("change_admin"), req.ID, resp.Admin))
	if err != nil {
		ah.Logger.Error("error replying to message", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(ah.Logger, m, botmsg)
}

func (ah *AdminHandlers) AllFlowers(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(ah.Logger, m, botmsg)
		return
	}

	var resp usersdata.GetAllFlowerTypesResp
	err := communication.MakeAdminHTTPReq(cfg.GetAllFlowerTypesURL, nil, &resp)
	if err != nil {
		ah.Logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}

	if resp.Err != "" {
		ah.Logger.Error("got error from resp", zap.Any("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}

	var text string
	for _, v := range resp.Result {
		text += fmt.Sprintf("%v:%v - %v\n", v.ID, v.Icon, v.Name)
	}
	text += fmt.Sprintf("len %v", len(resp.Result))
	botmsg, err := ah.Bot.Reply(m, text)
	if err != nil {
		ah.Logger.Error("error replying to message", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(ah.Logger, m, botmsg)
}

func (ah *AdminHandlers) RemoveFlower(m *telebot.Message) {
	if admin, err := ah.CheckAdmin(m.Sender.ID); !admin || err != nil {
		botmsg, _ := ah.Bot.Reply(m, localization.GetLoc("not_admin"))
		communication.UpdateUser(ah.Logger, m, botmsg)
		return
	}

	splitted := strings.Split(m.Text, " ")
	if len(splitted) != 2 {
		return
	}
	id, err := strconv.Atoi(splitted[1])
	if err != nil {
		ah.Logger.Error("error converting id", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.Bot.Reply(m, "error get id, need /removeFlower <id>")
		return
	}

	req := usersdata.RemoveFlowerReq{ID: types.Uint64(id)}
	var resp usersdata.RemoveFlowerResp
	err = communication.MakeAdminHTTPReq(cfg.RemoveFlowerURL, req, &resp)
	if err != nil {
		ah.Logger.Error("Error making admin request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.Bot.Reply(m, localization.GetLoc("error"))
		return
	}

	if !resp.OK {
		ah.Logger.Error("got error from resp", zap.Any("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		ah.Bot.Reply(m, localization.GetLoc("error"))
		return
	}
	botmsg, err := ah.Bot.Reply(m, "ok")
	if err != nil {
		ah.Logger.Error("error replying to message", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		return
	}
	go communication.UpdateUser(ah.Logger, m, botmsg)
}

func (ah *AdminHandlers) CheckAdmin(id int) (bool, error) {
	if id == Cfg.NeMoksID {
		return true, nil
	}
	req := usersdata.IsAdminReq{ID: id}
	var resp usersdata.IsAdminResp
	err := communication.MakeAdminHTTPReq(cfg.IsAdminURL, req, &resp)
	if err != nil {
		ah.Logger.Error("Error making admin request to user", zap.Error(err), zap.Int("user_id", id))
		return false, err
	}

	if resp.Err != "" {
		ah.Logger.Error("got error from resp", zap.Any("error", resp.Err), zap.Int("user_id", id))
		return false, fmt.Errorf(resp.Err)
	}
	return resp.Result, nil
}
