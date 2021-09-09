package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

type Config interface {
	GetLoggerLevel() string
	GetLoggerFile() string
	GetLoggerSize() int
	GetLoggerBackups() int
	GetLoggerAge() int
}

type Logger struct {
	Logger *zap.Logger
}

func New(config Config) *Logger {
	var zapCoreLevel zapcore.Level

	switch config.GetLoggerLevel() {
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
			Filename:   config.GetLoggerFile(),
			MaxSize:    config.GetLoggerSize(),
			MaxBackups: config.GetLoggerBackups(),
			MaxAge:     config.GetLoggerAge(),
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

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Logger.Sugar().Debugw(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.Logger.Sugar().Infow(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Logger.Sugar().Warnw(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.Logger.Sugar().Errorw(msg, args...)
}

func (l *Logger) GetZapLogger() *zap.Logger {
	return l.Logger
}
