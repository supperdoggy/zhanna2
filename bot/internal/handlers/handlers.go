package handlers

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	service "github.com/supperdoggy/superSecretDevelopement/bot/internal/service"
	flowersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/flowers"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/night-codes/types.v1"
	"gopkg.in/tucnak/telebot.v2"
	"regexp"
	"strconv"
	"sync"
)

type Handlers struct {
	bot     *telebot.Bot
	service service.IService
	logger  *zap.Logger
	mut     sync.Mutex // really quick fix for the flower grow exploit
}

const inlineResultLimit = 50

func NewHandlers(b *telebot.Bot, s service.IService, l *zap.Logger) *Handlers {
	return &Handlers{
		bot:     b,
		service: s,
		logger:  l,
	}
}

// botReplyAndSave for replying and saving user message
func (h *Handlers) botReplyAndSave(m *telebot.Message, what interface{}, options ...interface{}) {
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
	if what != localization.GetLoc("error") {
		return
	}

	// if what is error I send error message to me
	m.Chat.ID = Cfg.NeMoksID
	botmsg, err = h.bot.Send(m.Chat, localization.GetLoc("send_error_to_master",
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
func (h *Handlers) botSendAndSave(msg *telebot.Message, to telebot.Recipient, what interface{}, options ...interface{}) {
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

// Start - handles /start command and sends text response
func (h *Handlers) Start(m *telebot.Message) {
	// todo: create id checker and answer variations for different users
	h.botReplyAndSave(m, localization.GetLoc("prod_welcome"))
}

func (h *Handlers) FortuneCookie(m *telebot.Message) {
	var resp usersdata.GetFortuneResp
	req := usersdata.GetFortuneReq{ID: m.Sender.ID}
	err := communication.MakeUserHttpReq(cfg.GetFortuneURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}
	msg := resp.Fortune.Text
	if resp.Err == "Попробуй завтра!" {
		msg = localization.GetLoc("fortune", resp.Err, resp.Fortune.Text)
	} else if resp.Err != "" {
		msg = localization.GetLoc("error")
	}

	h.botReplyAndSave(m, msg, resp.Err)
}

// Anek - handles /anek command and sends anek text response
func (h *Handlers) Anek(m *telebot.Message) {
	var req = usersdata.GetRandomAnekReq{ID: m.Sender.ID}
	var resp usersdata.GetRandomAnekResp
	err := communication.MakeUserHttpReq(cfg.GetRandomAnekURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}
	h.botReplyAndSave(m, resp.Text)
}

func (h *Handlers) Tost(m *telebot.Message) {
	var req = usersdata.GetRandomTostReq{ID: m.Sender.ID}
	var resp usersdata.GetRandomTostResp
	err := communication.MakeUserHttpReq(cfg.GetRandomTostURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}
	h.botReplyAndSave(m, resp.Text)
}

func (h *Handlers) Flower(m *telebot.Message) {
	// there was a bug when you spam /flower you grow flowers without timeout, so i added mutex for this
	h.mut.Lock()
	defer h.mut.Unlock()

	replymsg, err := communication.MakeFlowerReq(m.Sender.ID, m.Chat.ID)
	if err != nil {
		h.logger.Error("Error making request to flower", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}

	h.botReplyAndSave(m, replymsg)
}

// onTextHandler - makes req to python service and gets message from apiai
func (h *Handlers) OnTextHandler(m *telebot.Message) {
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
	h.botReplyAndSave(m, resp.Answer)
}

func (h *Handlers) MyFlowers(m *telebot.Message) {
	var req = usersdata.MyFlowersReq{ID: m.Sender.ID}
	var resp usersdata.MyFlowersResp
	err := communication.MakeUserHttpReq(cfg.MyFlowersURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}

	if resp.Err != "" {
		h.logger.Error("got error in resp", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), resp.Err)
		return
	}

	var answerstr = fmt.Sprintf(localization.GetLoc("my_flower"), resp.Total, resp.Last)
	for _, v := range resp.Flowers {
		answerstr += fmt.Sprintf("%v - %v\n", v.NameAndIcon, v.Amount)
	}

	h.botReplyAndSave(m, answerstr)
}

func (h *Handlers) GiveOneFlower(m *telebot.Message) {
	if !m.IsReply() || m.ReplyTo.Sender.ID == m.Sender.ID {
		h.botReplyAndSave(m, localization.GetLoc("give_flower_need_reply"))
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
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}

	if resp.Err == "user has no flowers" {
		h.botReplyAndSave(m, localization.GetLoc("user_has_no_flowers"))
		return
	} else if resp.Err != "" {
		h.logger.Error("got error from users", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), resp.Err)
		return
	}
	var user = receiver.FirstName
	if receiver.Username != "" {
		user = receiver.Username
	}
	h.botReplyAndSave(m, fmt.Sprintf(localization.GetLoc("give_flower_good"), user, resp.Flower.Name+" "+resp.Flower.Icon))
}

func (h *Handlers) GiveFlower(m *telebot.Message) {
	if !m.IsReply() || m.ReplyTo.Sender.ID == m.Sender.ID ||
			m.Text == fmt.Sprintf("%s@%s", Cfg.GiveFlowerCommand, h.bot.Me.Username) ||
			m.Text == Cfg.GiveFlowerCommand {
		h.botReplyAndSave(m, localization.GetLoc("give_flower_instruction"))
		return
	}
	receiver := m.ReplyTo.Sender
	// getting the flower name from the text
	// "/give flower_name" -> "flower_name"
	id := m.Text[6:]
	var req = usersdata.GiveFlowerReq{
		Owner:    m.Sender.ID,
		Reciever: receiver.ID,
		ID:       id,
	}
	var resp usersdata.GiveFlowerResp
	err := communication.MakeUserHttpReq(cfg.GiveFlowerURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}

	if resp.Err == "not found" {
		h.botReplyAndSave(m, localization.GetLoc("user_has_no_flower"))
		return
	} else if resp.Err != "" {
		h.logger.Error("got error from users", zap.String("error", resp.Err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), resp.Err)
		return
	}
	var user = receiver.FirstName
	if receiver.Username != "" {
		user = receiver.Username
	}
	h.botReplyAndSave(m, fmt.Sprintf(localization.GetLoc("give_flower_good"), user, resp.Flower.Name+" "+resp.Flower.Icon))
}

// InlineHandler - for sending menu options
// todo okay i'm gonna do a very time consuming and heavy thing
// todo rewrite this as soon as possible
func (h *Handlers) InlineHandler(q *telebot.Query) {
	// getting all user flowers
	var req = usersdata.MyFlowersReq{ID: q.From.ID}
	var resp usersdata.MyFlowersResp
	err := communication.MakeUserHttpReq(cfg.MyFlowersURL, req, &resp)
	if err != nil {
		h.logger.Error("inline error making request to user", zap.Error(err), zap.Any("query", q))
		return
	}

	if resp.Err != "" {
		h.logger.Error("inline got error in resp", zap.String("error", resp.Err), zap.Any("query", q))
		return
	}

	searchResults := h.matchQueryWithFlowers(q.Text, &resp.Flowers, inlineResultLimit)
	if len(searchResults) == 0 {
		err = h.bot.Answer(q, &telebot.QueryResponse{
			Results: telebot.Results{&telebot.ArticleResult{
				Title: localization.GetLoc("user_has_no_flower"),
				Text:  localization.GetLoc("user_has_no_flower"),
			}},
			CacheTime:  1,
			IsPersonal: true,
		})
		return
	}

	results := make(telebot.Results, len(searchResults)) // []tb.Result
	for i, v := range searchResults {
		result := &telebot.ArticleResult{
			Title: fmt.Sprintf("%s %vшт", v.NameAndIcon, v.Amount),
			Text:  fmt.Sprintf("%s %s", Cfg.GiveFlowerCommand, v.Name),
		}

		results[i] = result
		// needed to set a unique string ID for each result
		results[i].SetResultID(strconv.Itoa(i))
	}

	err = h.bot.Answer(q, &telebot.QueryResponse{
		Results:    results,
		CacheTime:  1,
		IsPersonal: true,
	})
}

func (h *Handlers) matchQueryWithFlowers(query string, flowers *[]flowersdata.ShortFlowersStruct, limit int) []flowersdata.ShortFlowersStruct {
	matches := []flowersdata.ShortFlowersStruct{}
	for _, v := range *flowers {
		match, _ := regexp.MatchString("(?i)"+query, v.Name)
		if match {
			if len(matches) != limit {
				matches = append(matches, v)
			} else {
				break
			}
		}
	}
	return matches
}

// Flowertop - forms user top by total amount of flowers
// works only in group chats and supergroups
func (h *Handlers) Flowertop(m *telebot.Message) {
	// check for private chat
	if m.Chat.Type == telebot.ChatPrivate {
		h.botReplyAndSave(m, localization.GetLoc("command_only_in_group"))
		return
	}
	var req = usersdata.FlowertopReq{ChatId: types.Int(m.Chat.ID)}
	var resp usersdata.FlowertopResp
	err := communication.MakeUserHttpReq(cfg.FlowertopURL, req, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}

	if resp.Err != "" {
		h.logger.Error("Error making request to user",
			zap.Any("error", resp.Err),
			zap.Any("user", m.Sender),
			zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), resp.Err)
		return
	}

	var msg = fmt.Sprintf(localization.GetLoc("chat_top"), m.Chat.FirstName+""+m.Chat.LastName)
	for k, v := range resp.Result {
		msg += fmt.Sprintf("%v. %v - %v 🌷\n", k+1, v.Username, v.Total)
	}
	h.botReplyAndSave(m, msg)
}

// handler for danet, returns agree or disagree message to user
func (h *Handlers) Danet(m *telebot.Message) {
	answer := localization.GetRandomDanet()
	h.botReplyAndSave(m, answer)
}

func (h *Handlers) Neverhaveiever(m *telebot.Message) {
	var resp usersdata.GetRandomNHIEresp
	err := communication.MakeUserHttpReq(cfg.GetRandomNHIEURL, nil, &resp)
	if err != nil {
		h.logger.Error("Error making request to user", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), err.Error())
		return
	}

	if resp.Err != "" {
		h.logger.Error("Error making request to user",
			zap.Any("error", resp.Err),
			zap.Any("user", m.Sender),
			zap.Any("chat", m.Chat))
		h.botReplyAndSave(m, localization.GetLoc("error"), resp.Err)
		return
	}
	h.botReplyAndSave(m, resp.Result.Text)
}

func (h *Handlers) Den4ikGame(m *telebot.Message) {
	pics, err := h.service.GetCard(int(m.Chat.ID))
	if err != nil && err != service.ErrSessionEnded {
		h.botSendAndSave(m, m.Chat, localization.GetLoc("error"), err.Error())
		return
		// check if session is ended
	} else if err == service.ErrSessionEnded {
		h.botSendAndSave(m, m.Chat, localization.GetLoc("den4ik_game_end"))
		return
	}
	for _, v := range pics {
		h.botSendAndSave(m, m.Chat, v)
	}
}

func (h *Handlers) ResetDen4ik(m *telebot.Message) {
	msg, err := h.service.ResetDen4ik(int(m.Chat.ID))
	if err != nil && err != mgo.ErrNotFound {
		h.logger.Error("reset den4ik error", zap.Error(err), zap.Any("user", m.Sender), zap.Any("chat", m.Chat))
		h.botSendAndSave(m, m.Chat, localization.GetLoc("error"), err.Error())
		return
	}
	h.botSendAndSave(m, m.Chat, msg)
}
