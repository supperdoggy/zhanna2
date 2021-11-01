package localization

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type localization struct {
	m   map[string]map[string]string
	mut sync.Mutex
	danet map[string][]string
	mutdanet sync.Mutex
}

var (
	loc = localization{
		m: map[string]map[string]string{
			"ru": {
				"zhanna_has_flower": "Жанночки",
				"error":                   "что-то пошло по пизде сори, напиши пж этому крутому челу -> @supperdoggy",
				"send_error_to_master":    "hey daddy, тут у чела @%s траблы, глянь пж\n\nuser: %+v \n\nchat: %+v \n\noptions: %+v",
				"command_only_in_group":   "команда доступна только в груповом чате",
				"give_flower_good":        "ты успешно подарил цветок! \nв коллекции %v теперь есть %v",
				"give_flower_need_reply":  "Тебе нужно ответить на сообщение человека которому ты хочешь подарить цветок!",
				"user_has_no_flowers":     "У тебя пока что нету цветочков чтобы подарить кому-то...",
				"user_has_no_flower":      "у автора нет этого цветочка соре",
				"give_flower_instruction": "Тебе нужно ответить на сообщение человека, которому ты хочешь подарить цветок, а затем написать @Zhanno4kabot и выбрать в меню цветочек :)",
				"my_flower":               "Вот твои цветочки!\nУ тебя уже %v 🌷 %v 🌱\n\n",
				"add_flower":              "неправильный формат, надо иконка-название-категория",
				"chat_top":                "Вот топ чатика: %v\n\n",
				"dev_welcome":             "Привет, я киберзомби жанночка, скоро я захвачу мир, если ноут не закроют...",
				"prod_welcome":            "Приветик, я жанночка, давай знакомиться?",
				"yes":                     "да",
				"no":                      "нет",
				"idk":                     "хз",
				"already_grew_flowers":    "Ты уже сегодня поливал цветочки!\nПопробуй позже",
				"flower_grew":             "Поздравляю! Твой %v вырос!",
				"flower_grew_not_fully":   "Твой цветок вырос на %v цветочковых единиц!\nДополнительный прирост: %v",
				"flower_already_have":     "\nСейчас у тебя %v🌷 и %v🌱",
				"not_admin":               "Вы не администратор",
				"need_reply":              "Ответь на сообщение кого надо сделать админом",
				"flower_died":             "Штош, у меня плохие новости, твой цветок умер",
				"fortune":                 "%v\n\n%v",
				"admin_help": "/admin - set/unset admin\n" +
					"/addFlower - add new flower type\n" +
					"/removeFlower - remove flower type\n" +
					"/allFlowers - returns flower types list\n",
				"done":                  "Дело сделано!",
				"change_admin":          "Пользователь %v admin: %v\n",
				"9_card":                "Игрок, вытянувший карту, загадывает тему (футбол/напитки/финансы, что угодно), далее все по очереди называют слова связанные с этой темой, тот, кому нечего сказать или думает больше 5 сек. - пьет",
				"10_card":               "Быстро хлопни в ладони! Кто последний хлопнет, тот пьет",
				"j_card":                "Это твое право на перерыв (2 мин). Пока тебя нет игра продолжается",
				"q_card":                "Теперь ты изгой. Другим игрокам нельзя с тобой говорить, если они это сделают или ответят на твой вопрос - они пьют. Можешь провоцировать :)\nДама действует до тех пор пока не выпала следующая",
				"k_card":                "Игрок, вытянувший карту, придумывает движение, которое нужно исполнить перед тем, как пить дальше. С каждым новым королём добавляется ещё одно движение — вплоть до четырёх в строго установленном порядке. Тот, кто забыл сделать движение — пьёт ещё раз.",
				"a_card":                "Игрок, который вытянул туз выбирает человека, который сейчас будет пить",
				"den4ik_game_end":       "Карты закончились. Сыграем еще раз?",
				"reset_ok":              "Ты успешно сбросил сессию игры денчика! Сыграем сначала?",
				"remove_flower_need_id": "ошибка получения айди, надо использовать шаблон /removeFlower <id>",
			},
			"uk": {
				"zhanna_has_flower": "Жанночки",
				"error":                   "щось пішло по пизді сорі, напиши буласка цьому крутому чєлу -> @supperdoggy",
				"send_error_to_master":    "hey daddy, тут у чє0ла @%s трабли, подивись буласка\n\nuser: %+v \n\nchat: %+v \n\noptions: %+v",
				"command_only_in_group":   "команда доступна лише в груповому чаті",
				"give_flower_good":        "ти успішно подарував квітку! \nв колекції %v тепер є %v",
				"give_flower_need_reply":  "Тобі потрібно відповісти на повідомлення людини, якій ти хочеш подарувати квітку!",
				"user_has_no_flowers":     "У тебе поки що немає квіточок, щоб дарувати комусь…",
				"user_has_no_flower":      "у автора немає цієї квіточки сорі",
				"give_flower_instruction": "Тобі потрібно відповісти на повідомлення людини, якій ти хочеш подарувати квітку, а потім написати @Zhanno4kabot и обрати квіточку в меню :)",
				"my_flower":               "Ось твої квіточки! \nУ тебе вже %v 🌷 %v 🌱\n\n",
				"add_flower":              "неправильний формат, потрібно іконка-назва-категорія",
				"chat_top":                "Ось топ чатіка: %v\n\n",
				"dev_welcome":             "Привіт, я кіберзомбі жанночка, скоро я захоплю світ, якщо ноут не закриють...",
				"prod_welcome":            "Привітики, я Жанночка, нумо знайомитися?",
				"yes":                     "так",
				"no":                      "ні",
				"idk":                     "хз",
				"already_grew_flowers":    "Ти уже сьогодні поливав квіточки!\nСпробуй пізніше",
				"flower_grew":             "Вітаю! Твоя %v виросла!",
				"flower_grew_not_fully":   "Твоя квітка виросла на %v квіточкових одиниць!\nДодатковий приріст: %v",
				"flower_already_have":     "\nЗараз у тебе %v🌷 і %v🌱",
				"not_admin":               "Ви не адміністратор",
				"need_reply":              "Відповідай на повідомлення того, кого потрібно зробити адміном",
				"flower_died":             "Штош, у мене погані новини, твоя квітка померла",
				"fortune":                 "%v\n\n%v",
				"admin_help": "/admin - set/unset admin\n" +
					"/addFlower - add new flower type\n" +
					"/removeFlower - remove flower type\n" +
					"/allFlowers - returns flower types list\n",
				"done":                  "Справа зроблена!",
				"change_admin":          "Користувач %v admin: %v\n",
				"9_card":                "Гравець, що витягнув карту, загадує тему (футбол/напої/фінанси тощо), а далі усі по черзі називають слова, що пов’язані з цією темою. Той, кому немає що сказати або думає більше 5 секунд – п’є",
				"10_card":               "Швидко хлопни долонями! Хто хлопне останнім, той п’є",
				"j_card":                "Це твоє право на перерву (2 хв). Поки тебе немає, гра продовжується",
				"q_card":                "Тепер ти вигнанець. Іншим гравцям не можна з тобою говорити, але якщо вони це зроблять або дадуть відповідь на твоє питання – вони п’ють. Можеш провокувати :)\nДама діє до тих пір, поки не випала наступна",
				"k_card":                "Гравець, що витягнув карту, придумує рух, який необхідно виконати перед тим, як пити далі. З кожним новим королем додається ще один рух – аж до чотирьох у строго встановленому порядку. Той, хто забув виконати рух – п’є ще раз.",
				"a_card":                "Гравець, що витягнув туз, обирає людину, яка зараз буде пити",
				"den4ik_game_end":       "Карти скінчились. Зіграємо ще раз?",
				"reset_ok":              "Ти успішно скинув сесію гри денчіка. Зіграємо спочатку?",
				"remove_flower_need_id": "помилка отримання айді, потрібно використовувати шаблон /removeFlower <id>",
			},
		},
		danet: map[string][]string{
			"ru": {"Мама сказала ", "Видимо ", "Точно ", "Я сказала ", "Походу ", "Мне мама сказала что ", "Надеюсь что ", "Звезды сказали "},
			"uk": {"Мама сказала ", "Схоже що ", "Точно ", "Я сказала ", "Походу ", "Мені мама сказала що ", "Надіюся що ", "Зорі сказали "},
		},
	}
)

const NoLocalizationError = "No Localization found"

// returns localization by lang and key
func GetLoc(key string, lang string, args ...interface{}) string {
	// default language is ru
	if _, ok := loc.m[lang]; !ok {
		lang = "ru"
	}

	loc.mut.Lock()
	defer loc.mut.Unlock()
	val, ok := loc.m[lang][key]
	if !ok {
		val = NoLocalizationError
	}
	if len(args) != 0 {
		val = fmt.Sprintf(val, args...)
	}
	return val
}

// returns string with random agree or disagree
func GetRandomDanet(lang string) string {
	// default language is ru
	if lang != "uk" && lang != "ru" {
		lang = "ru"
	}

	loc.mutdanet.Lock()
	defer loc.mutdanet.Unlock()
	rand.Seed(time.Now().UnixNano())
	agree := rand.Intn(2) == 0
	if agree {
		return loc.danet[lang][rand.Intn(len(loc.danet))] + GetLoc("yes", lang)
	}
	return loc.danet[lang][rand.Intn(len(loc.danet))] + GetLoc("no", lang)
}
