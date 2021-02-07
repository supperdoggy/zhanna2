package main

import "sync"

type localization struct {
	m   map[string]string
	mut sync.Mutex
}

var loc = localization{
	m: map[string]string{
		"error":                  "что-то пошло по пизде сори",
		"command_only_in_group":  "комманда доступна только в груповом чате",
		"give_flower_good":       "ты успешно подарил цветочек",
		"give_flower_need_reply": "Тебе нужно ответить на сообщение человека которому ты хочешь подарить цветок!",
		"my_flower":              "Вот твои цветочки!\nУ тебя уже %v 🌷 %v 🌱\n\n",
		"add_flower":             "неправильный формат, надо иконка-название-категория",
		"chat_top":               "Вот топ чатика: %v\n\n",
		"dev_welcome":            "Привет, я пока что очень сырая, будь нежен со мной...",
	},
}

func getLoc(key string) string {
	loc.mut.Lock()
	defer loc.mut.Unlock()
	val := loc.m[key]
	return val
}
