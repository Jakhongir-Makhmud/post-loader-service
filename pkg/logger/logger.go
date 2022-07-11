package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Int ...
	Int = zap.Int
	// String ...
	String = zap.String
	// Error ...
	Error = zap.Error
	// Bool ...
	Bool = zap.Bool

	// Any ...
	Any = zap.Any
)

type Field = zapcore.Field

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Fatal(string, ...Field)
}

type loggerImpl struct {
	zap *zap.Logger
}

var (
	customTimeFormat string
)

func (l *loggerImpl) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

func (l *loggerImpl) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

func (l *loggerImpl) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

func (l *loggerImpl) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

func (l *loggerImpl) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, fields...)
}

func New(level string, namespace string) Logger {
	if level == "" {
		level = LevelInfo
	}

	logger := loggerImpl{
		zap: newZapLogger(level, time.RFC3339),
	}

	logger.zap = logger.zap.Named(namespace)

	zap.RedirectStdLog(logger.zap)

	return &logger
}

// GetNamed ...
func GetNamed(l Logger, name string) Logger {
	switch v := l.(type) {
	case *loggerImpl:
		v.zap.Named(name)
		return v
	default:
		l.Info("logger.GetNamed: Invalid logger type")
		return l
	}
}

// WithFields ...
func WithFields(l Logger, fields ...Field) Logger {
	switch v := l.(type) {
	case *loggerImpl:
		return &loggerImpl{
			zap: v.zap.With(fields...),
		}
	default:
		l.Info("logger.WithFields: Invalid logger type")
		return l
	}
}

// Cleanup ...
func Cleanup(l Logger) error {
	switch v := l.(type) {
	case *loggerImpl:
		return v.zap.Sync()
	default:
		l.Info("logger.Cleanup: Invalid logger type")
		return nil
	}
}
