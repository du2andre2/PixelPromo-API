package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *Config) Logger {
	lgCfg := zap.NewProductionConfig()
	if cfg.Env == Local {
		lgCfg = zap.NewDevelopmentConfig()
	}

	lg, err := lgCfg.Build()
	if err != nil {
		panic(err)
	}
	return logger{
		lg: lg,
	}
}

type level int8

const (
	levelDebug level = iota - 1
	levelInfo
	levelWarn
	levelError
	levelPanic
	levelFatal
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Flush() error
}

type Field struct {
	Key string
	Val interface{}
}

func F(key string, val interface{}) Field {
	return Field{Key: key, Val: val}
}

type logger struct {
	lg *zap.Logger
}

func (l logger) Debug(msg string, fields ...Field) {
	l.Log(levelDebug, msg, fields...)
}

func (l logger) Info(msg string, fields ...Field) {
	l.Log(levelInfo, msg, fields...)
}

func (l logger) Warn(msg string, fields ...Field) {
	l.Log(levelWarn, msg, fields...)
}

func (l logger) Error(msg string, fields ...Field) {
	l.Log(levelError, msg, fields...)
}

func (l logger) Panic(msg string, fields ...Field) {
	l.Log(levelPanic, msg, fields...)
}

func (l logger) Fatal(msg string, fields ...Field) {
	l.Log(levelFatal, msg, fields...)
}

func (l logger) Flush() error {
	return l.lg.Sync()
}

func (l logger) Log(lvl level, msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Val)
	}
	l.lg.WithOptions(zap.AddCallerSkip(2)).Log(zapcore.Level(lvl), msg, zapFields...)
}
