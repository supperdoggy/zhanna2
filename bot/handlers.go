package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/tucnak/telebot.v2"
)

func testMessage(m *telebot.Message) {
	resp := fmt.Sprintf("%v", m.IsReply())

	b, _ := bot.Reply(m, resp)
	log.Println(b.Sender.ID)
}

// start() - handles /start command and sends text response
// todo below
func start(m *telebot.Message) {
	var response string
	// todo: create id checker and answer variations for different users
	response = "Привет, я пока что очень сырая, будь нежен со мной..."
	botmsg, err := bot.Reply(m, response)
	if err != nil {
		log.Println("handlers.go -> start() -> error:", err.Error(), ", user id:", m.Sender.ID)
		return
	}
	go UpdateUser(m, botmsg)
}

// TODO:
func fortuneCookie(m *telebot.Message) {
	var resp struct {
		FortuneCookie
		Err string `json:"err"`
	}
	data := obj{"id": m.Sender.ID}
	readyData, err := json.Marshal(data)
	if err != nil {
		log.Println("error unmarshalling")
		return
	}
	r, err := MakeHttpReq(userUrl+"/getFortune", "POST", readyData)
	if err != nil {
		log.Println("error making request")
		return
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		log.Println("error unmarshalling")
		return
	}
	msg := resp.Text
	if resp.Err != "" {
		msg = resp.Err
	}

	botmsg, err := bot.Reply(m, msg)
	if err != nil {
		log.Println("error sending answer, FortuneCookie:", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}

// anek() - handles /anek command and sends anek text response
func anek(m *telebot.Message) {
	anekAnswer, err := MakeRandomAnekHttpReq(m.Sender.ID)
	if err != nil {
		log.Println("handlers.go -> anek() -> make req error:", err.Error())
		return
	}
	botmsg, err := bot.Reply(m, anekAnswer.Text)
	if err != nil {
		log.Println("handlers.go -> anek() -> reply error:", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}

func tost(m *telebot.Message) {
	answerTost, err := MakeRandomTostHttpReq(m.Sender.ID)
	if err != nil {
		log.Println("handlers.go -> tost() -> make req error:", err.Error())
		return
	}
	botmsg, err := bot.Reply(m, answerTost.Text)
	if err != nil {
		log.Println("handlers.go -> tost() -> reply error:", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}

func addFlower(m *telebot.Message) {
	text := split(m.Text[11:], "-")
	if len(text) != 3 {
		bmsg, _ := bot.Reply(m, "wrong format, need text-text-text")
		go UpdateUser(m, bmsg)
		return
	}
	data := obj{"icon": text[0], "name": text[1], "type": text[2]}
	_, err := MakeUserHttpReq("addFlower", data)
	if err != nil {
		log.Println("handlers.go -> addFlower() -> MakeUserHttpReq error:", err.Error())
		botmsg, _ := bot.Reply(m, "communication error")
		go UpdateUser(m, botmsg)
		return
	}
	botmsg, _ := bot.Reply(m, "Done!")
	go UpdateUser(m, botmsg)
}

func flower(m *telebot.Message) {
	resp, err := MakeFlowerReq(m.Sender.ID)
	if err != nil {
		log.Println("handlers.go -> flower() -> MakeFlowerReq() error", err.Error(), m.Sender.ID)
		_, _ = bot.Reply(m, "error occured, contact owner")
		return
	}
	botmsg, err := bot.Reply(m, resp)
	if err != nil {
		log.Println("handlers.go -> flower() -> bot.Reply() error", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}

func onTextHandler(m *telebot.Message) {
	answer, err := MakeUserHttpReq("getAnswer", obj{"id": m.Sender.ID, "text": m.Text})
	if err != nil {
		log.Println("onTextHandler() -> req error:", err.Error())
		bot.Reply(m, "Error getting answer")
		return
	}

	var resp struct {
		Answer string `json:"answer"`
		Err    string `json:"err"`
	}
	if err := json.Unmarshal(answer, &resp); err != nil {
		log.Println("onTextHandler() -> Unmarshal error:", err.Error())
		bot.Reply(m, "Error unmarhsal")
		return
	}
	if resp.Err != "" {
		log.Println("onTextHandler() -> got error in response:", resp.Err)
		bot.Reply(m, resp.Err)
		return
	}
	botmsg, _ := bot.Reply(m, resp.Answer)
	go UpdateUser(m, botmsg)
}

func myflowers(m *telebot.Message) {
	answer, err := MakeUserHttpReq("myflowers", obj{"id": m.Sender.ID})
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
		bot.Reply(m, resp.Err)
		return
	}

	var answerstr string = fmt.Sprintf("Вот твои цветочки!\nУ тебя уже %v 🌷 %v 🌱\n\n", resp.Total, resp.Last)
	for k, v := range resp.Flowers {
		answerstr += fmt.Sprintf("%v - %v\n", k, v)
	}
	botmsg, _ := bot.Reply(m, answerstr)
	go UpdateUser(m, botmsg)
}

func giveOneFlower(m *telebot.Message) {
	if !m.IsReply() {
		b, _ := bot.Reply(m, "Тебе нужно ответить на сообщение человека которому ты хочешь подарить цветок!")
		UpdateUser(m, b)
		return
	}

	data := obj{"last": true, "owner": m.Sender.ID, "reciever": m.ReplyTo.Sender.ID}
	answer, err := MakeUserHttpReq("give", data)
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
	b, _ := bot.Reply(m, "Ты успешно подарил цветок!")
	go UpdateUser(m, b)
}

// forms user top by total amount of flowers
// works only in group chats and supergroups
func flowertop(m *telebot.Message) {
	// check for private chat
	if m.Chat.Type == telebot.ChatPrivate {
		botmsg, _ := bot.Reply(m, "Функция доступна только в груповых чатах")
		UpdateUser(m, botmsg)
		return
	}
	answer, err := MakeUserHttpReq("flowertop", obj{"chatid": m.Chat.ID})
	if err != nil {
		log.Printf("handlers.go -> flowertop() -> MakeUserHttpReq('flowertop') error: %v, chatid: %v\n", err.Error(), m.Chat.ID)
		botmsg, _ := bot.Reply(m, "Что-то пошло по пизде сори")
		UpdateUser(m, botmsg)
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
		log.Printf("handlers.go -> flowertop() -> Unmarshal error:%v, body: %v\n", err.Error(), string(answer))
		botmsg, _ := bot.Reply(m, "Что-то пошло по пизде, а именно анмаршал(напиши максу он скажет что не так])")
		UpdateUser(m, botmsg)
		return
	}
	var msg string = fmt.Sprintf("Вот топ чатика: %v\n\n", m.Chat.FirstName+""+m.Chat.LastName)
	for k, v := range resp.Top {
		msg += fmt.Sprintf("%v. %v - %v 🌷\n", k+1, v.Username, v.Total)
	}
	botmsg, _ := bot.Reply(m, msg)
	UpdateUser(m, botmsg)
}
