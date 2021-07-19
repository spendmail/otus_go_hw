package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

type Logger struct {
	Logger *zap.Logger
}

func New(config internalconfig.LoggerConf) *Logger {
	var zapCoreLevel zapcore.Level

	switch config.Level {
	case DEBUG:
		zapCoreLevel = zap.DebugLevel
	case INFO:
		zapCoreLevel = zap.InfoLevel
	case WARN:
		zapCoreLevel = zap.WarnLevel
	case ERROR:
		zapCoreLevel = zap.ErrorLevel
	default:
		zapCoreLevel = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.File,
			MaxSize:    config.Size,
			MaxBackups: config.Backups,
			MaxAge:     config.Age,
		}),
		zapCoreLevel,
	)
	logger := zap.New(core)

	return &Logger{
		Logger: logger,
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (l *Logger) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l *Logger) Info(msg string) {
	l.Logger.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.Logger.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.Logger.Error(msg)
}
