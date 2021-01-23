package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/tucnak/telebot.v2"
)

// start() - handles /start command and sends text response
// todo below
func start(m *telebot.Message) {
	var response string
	// todo: create id checker and answer variations for different users
	response = "Привет, я пока что очень сырая, будь нежен со мной..."
	botmsg, err := bot.Reply(m, response)
	if err != nil {
		fmt.Println("handlers.go -> start() -> error:", err.Error(), ", user id:", m.Sender.ID)
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
		fmt.Println("error unmarshalling")
		return
	}
	r, err := MakeHttpReq(userUrl+"/getFortune", "POST", readyData)
	if err != nil {
		fmt.Println("error making request")
		return
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println("error unmarshalling")
		return
	}
	msg := resp.Text
	if resp.Err != "" {
		msg = resp.Err
	}

	botmsg, err := bot.Reply(m, msg)
	if err != nil {
		fmt.Println("error sending answer, FortuneCookie:", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}

// anek() - handles /anek command and sends anek text response
func anek(m *telebot.Message) {
	anekAnswer, err := MakeRandomAnekHttpReq(m.Sender.ID)
	if err != nil {
		fmt.Println("handlers.go -> anek() -> make req error:", err.Error())
		return
	}
	botmsg, err := bot.Reply(m, anekAnswer.Text)
	if err != nil {
		fmt.Println("handlers.go -> anek() -> reply error:", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}

func tost(m *telebot.Message) {
	answerTost, err := MakeRandomTostHttpReq(m.Sender.ID)
	if err != nil {
		fmt.Println("handlers.go -> tost() -> make req error:", err.Error())
		return
	}
	botmsg, err := bot.Reply(m, answerTost.Text)
	if err != nil {
		fmt.Println("handlers.go -> tost() -> reply error:", err.Error())
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
	marhshaled, err := json.Marshal(data)
	if err != nil {
		fmt.Println("handlers.go -> addFlower() -> marshal error:", err.Error())
		botmsg, _ := bot.Reply(m, "unmarshal error")
		go UpdateUser(m, botmsg)
		return
	}
	_, err = MakeUserHttpReq("addFlower", marhshaled)
	if err != nil {
		fmt.Println("handlers.go -> addFlower() -> MakeUserHttpReq error:", err.Error())
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
		fmt.Println("handlers.go -> flower() -> MakeFlowerReq() error", err.Error(), m.Sender.ID)
		_, _ = bot.Reply(m, "error occured, contact owner")
		return
	}
	botmsg, err := bot.Reply(m, resp)
	if err != nil {
		fmt.Println("handlers.go -> flower() -> bot.Reply() error", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}

func onTextHandler(m *telebot.Message) {
	data, err := json.Marshal(obj{"id": m.Sender.ID, "text": m.Text})
	if err != nil {
		fmt.Println("onTextHandler() -> Marshal error:", err.Error())
		bot.Reply(m, "Error getting answer")
		return
	}

	answer, err := MakeUserHttpReq("getAnswer", data)
	if err != nil {
		fmt.Println("onTextHandler() -> req error:", err.Error())
		bot.Reply(m, "Error getting answer")
		return
	}

	var resp struct {
		Answer string `json:"answer"`
		Err    string `json:"err"`
	}
	if err := json.Unmarshal(answer, &resp); err != nil {
		fmt.Println("onTextHandler() -> Unmarshal error:", err.Error())
		bot.Reply(m, "Error unmarhsal")
		return
	}
	if resp.Err != "" {
		fmt.Println("onTextHandler() -> got error in response:", resp.Err)
		bot.Reply(m, resp.Err)
		return
	}
	botmsg, _ := bot.Reply(m, resp.Answer)
	go UpdateUser(m, botmsg)
}

func myflowers(m *telebot.Message) {
	data, err := json.Marshal(obj{"id": m.Sender.ID})
	if err != nil {
		fmt.Println("myflowers() -> Marshal error:", err.Error())
		return
	}
	answer, err := MakeUserHttpReq("myflowers", data)
	if err != nil {
		fmt.Println("myflowers() -> MakeUserHttpReq(myflowers) err:", err.Error())
		return
	}
	var resp struct {
		Flowers map[string]int `json:"flowers"`
		Last    uint8          `json:"last"`
		Total   int            `json:"total"`
		Err     string         `json:"err"`
	}

	if err := json.Unmarshal(answer, &resp); err != nil {
		fmt.Println("myflowers() -> unmarshal error:", err.Error(), string(answer))
		return
	}

	if resp.Err != "" {
		fmt.Println("myflowers() -> got error resp from service:", resp.Err)
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
