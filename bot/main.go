package main

import (
	"fmt"
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

	fmt.Println("Bot running...")
	bot.Start()
}
