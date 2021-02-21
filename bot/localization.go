package main

import (
	"math/rand"
	"sync"
	"time"
)

type localization struct {
	m   map[string]string
	mut sync.Mutex
}

var (
	loc = localization{
		m: map[string]string{
			"error":                  "что-то пошло по пизде сори",
			"command_only_in_group":  "комманда доступна только в груповом чате",
			"give_flower_good":       "ты успешно подарил цветочек",
			"give_flower_need_reply": "Тебе нужно ответить на сообщение человека которому ты хочешь подарить цветок!",
			"my_flower":              "Вот твои цветочки!\nУ тебя уже %v 🌷 %v 🌱\n\n",
			"add_flower":             "неправильный формат, надо иконка-название-категория",
			"chat_top":               "Вот топ чатика: %v\n\n",
			"dev_welcome":            "Привет, я пока что очень сырая, будь нежен со мной...",
			"yes":                    "да",
			"no":                     "нет",
			"idk":                    "хз",
			"already_grew_flowers":   "Ты уже сегодня поливал цветочки!\nПопробуй позже",
			"flower_grew":            "Поздравляю! Твой %v вырос!",
			"flower_grew_not_fully":  "Твой цветок вырос на %v цветочковых единиц!",
			"flower_already_have":    "\nСейчас у тебя %v🌷 и %v🌱",
			"not_admin":              "Вы не администратор",
			"need_reply":             "Ответь на сообщение кого надо сделать админом",
			"flower_died":            "Штош, у меня плохие новости, твой цветок умер",
		},
	}
	// for danet
	danetVariations []string = []string{"Мама сказала ", "Видимо ", "Точно ", "Я сказала ", "Походу ", "Мне мама сказала что ", "Надеюсь что ", "Звезды сказали "}
)

// returns localization
func getLoc(key string) string {
	loc.mut.Lock()
	defer loc.mut.Unlock()
	val := loc.m[key]
	return val
}

// returns string with random agree or disagree
func getRandomDanet() string {
	rand.Seed(time.Now().UnixNano())
	agree := rand.Intn(2) == 0
	if agree {
		return danetVariations[rand.Intn(len(danetVariations))] + getLoc("yes")
	}
	return danetVariations[rand.Intn(len(danetVariations))] + getLoc("no")
}