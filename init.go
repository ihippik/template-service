package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ihippik/template-service/config"
)

func iniLogger(cfg config.LogCfg, version string) (*zap.Logger, error) {
	lCfg := zap.NewProductionConfig()
	lCfg.DisableCaller = !cfg.Caller
	lCfg.DisableStacktrace = !cfg.StackTrace
	lCfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	}

	switch cfg.Level {
	case "debug":
		lCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "warn":
		lCfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		lCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, err := lCfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.With(zap.Field{Key: "version", Type: zapcore.StringType, String: version}), nil
}

func initConn(cfg config.DBCfg) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.Conn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}
