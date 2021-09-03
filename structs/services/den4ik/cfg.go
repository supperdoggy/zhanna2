package den4ikcfg

import defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"

const (
	Port                   = ":9999"
	GetCardURL             = "/getCard"
	SessionReset           = "/resetSession"
	DBName                 = "Den4ik"
	GameSessionsCollection = "Sessions"
	URL                    = "http://localhost" + Port + "/" + defaultCfg.ApiV1
)
