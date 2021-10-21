package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"sync"

)

// config - Application config
type config struct {
	Token string `env:"BOT_TOKEN,required"`
}

var c config

func GetConfig(logger *zap.Logger) *config {
	ctx := context.Background()
	once := sync.Once{}
	once.Do(func() {
		if err := envconfig.Process(ctx, &c); err != nil {
			logger.Fatal("error processing config", zap.Error(err))
		}
	})
	return &c
}
