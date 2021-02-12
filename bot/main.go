package main

import (
	"fmt"

	"gopkg.in/tucnak/telebot.v2"
)

type obj map[string]interface{}

var (
	bot = initBot()
	// db itself
	DB = DbStruct{
		DbSession: connectToDB(),
	}
)

func main() {
	// connection to collections
	DB.CfgCollection = connectToCfgCollection()
	DB.MessagesCollection = connectToMessagesCollection()
	DB.PhrasesCollection = connectToPhrasesCollection()

	// handlers
	bot.Handle("/start", start)
	bot.Handle("/fortune", fortuneCookie)
	bot.Handle("/anek", anek)
	bot.Handle("/tost", tost)
	// bot.Handle("/addFlower", addFlower)
	bot.Handle("/flower", flower)
	bot.Handle(telebot.OnText, onTextHandler)
	bot.Handle("/myflowers", myflowers)
	bot.Handle("/giveoneflower", giveOneFlower)
	bot.Handle("/testMessage", testMessage)
	bot.Handle("/flowertop", flowertop)
	bot.Handle("/danet", danet)

	// admin handlers
	bot.Handle("/addFlower", addFlower)
	bot.Handle("/admin", admin)

	fmt.Println("Bot running...")
	bot.Start()
}
