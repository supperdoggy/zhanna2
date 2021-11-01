package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"sync"
)

// config - Application config
type config struct {
	Token                  string `env:"BOT_TOKEN,required"`
	ErrorAdminNotification bool   `env:"ERROR_ADMIN_NOTIFICATION"`
	IsProd bool `env:"IS_PROD"`
}

var c config

func GetConfig(logger *zap.Logger) *config {
	once := sync.Once{}
	once.Do(func() {
		ctx := context.Background()
		if err := envconfig.Process(ctx, &c); err != nil {
			logger.Fatal("error processing config", zap.Error(err))
		}
	})
	return &c
}
