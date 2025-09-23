package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var (
	Int64 = zap.Int64
	Error = zap.Error
)

type Config struct {
	Level            string   `mapstructure:"level"`
	OutputPaths      []string `mapstructure:"out_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

type Logger interface {
	Error(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Debug(msg string, fields ...zapcore.Field)
}

type logger struct {
	log *zap.Logger
}

func (l *logger) Error(msg string, fields ...zapcore.Field) {
	l.log.Error(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zapcore.Field) {
	l.log.Warn(msg, fields...)
}

func (l *logger) Info(msg string, fields ...zapcore.Field) {
	l.log.Info(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zapcore.Field) {
	l.log.Debug(msg, fields...)
}

func New(conf Config) (Logger, error) {
	logLevel, err := zap.ParseAtomicLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	cfg := zap.Config{
		Level:            logLevel,
		Encoding:         "json",
		Development:      false,
		OutputPaths:      conf.OutputPaths,
		ErrorOutputPaths: conf.ErrorOutputPaths,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}
	zapLogger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &logger{log: zapLogger}, nil
}

func NewDefault() Logger {
	logg, err := New(
		Config{
			Level:            "INFO",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		},
	)
	if err != nil {
		log.Fatal("default logger was failed. Happens something critical.")
	}
	logg.Info("a default logger was initialized.")
	return logg
}

func StringVal(msg string, val string) zap.Field {
	return zap.String(msg, val)
}

func IntVal(msg string, val int) zap.Field {
	return zap.Int(msg, val)
}

func ErrorVal(msg string, val error) zap.Field {
	return zap.NamedError(msg, val)
}
