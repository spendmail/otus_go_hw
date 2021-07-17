package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

type Logger struct {
	File    string
	Level   string
	Size    int
	Backups int
	Age     int
	Logger  *zap.Logger
}

func New(file, level string, size, backups, age int) *Logger {

	var zapCoreLevel zapcore.Level

	switch level {
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
			Filename:   file,
			MaxSize:    size,
			MaxBackups: backups,
			MaxAge:     age,
		}),
		zapCoreLevel,
	)
	logger := zap.New(core)

	return &Logger{
		File:    file,
		Level:   level,
		Size:    size,
		Backups: backups,
		Age:     age,
		Logger:  logger,
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
