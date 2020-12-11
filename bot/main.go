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
	fmt.Println("Handlers init start")

	bot.Handle("/start", start)
	bot.Handle("/fortuneCookie", fortuneCookie)
	bot.Handle("/anek", anek)

	fmt.Println("Bot running...")
	bot.Start()
}
