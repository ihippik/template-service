package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type (
	Config struct {
		ServerAddr string `env:"SERVER_ADDR,required"`
		DB         DBCfg  `env:",prefix=DB_"`
		Log        LogCfg `env:",prefix=LOG_"`
	}

	LogCfg struct {
		Level      string `env:"LEVEL,default=info"`
		Caller     bool   `env:"CALLER"`
		StackTrace bool   `env:"STACK_TRACE"`
	}

	DBCfg struct {
		Conn         string `env:"CONN,required"`
		MaxOpenConns int    `env:"MAX_OPEN_CONNS, default=10"`
		MaxIdleConns int    `env:"MAX_IDLE_CONNS, default=10"`
	}
)

func New(ctx context.Context) (*Config, error) {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
