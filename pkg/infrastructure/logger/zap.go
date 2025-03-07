package logger

import (
	"go-clean-api/pkg"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	inner *zap.Logger
}

// NewZapLogger creates a new custom Zap logger.
func NewZapLogger(config pkg.Config) (*ZapLogger, error) {
	// Logs outputs
	outputs, err := getLoggerOutputs(config.Log.Outputs, config.AppName, config.Log.Path)
	if err != nil {
		return nil, err
	}

	// Level
	level := getZapLoggerLevel(config.Log.Level, config.AppEnv)

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      outputs,
		ErrorOutputPaths: outputs,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.RFC3339TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	defaultLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	logger, err := cfg.Build()
	if err != nil {
		return &ZapLogger{defaultLogger}, nil
	}

	return &ZapLogger{logger}, nil
}

// FromFields converts fields to zap.Field.
func (l *ZapLogger) FromFields(fields Fields) any {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch f.Type {
		case "int":
			zapFields[i] = zap.Int(f.Key, f.Value.(int))
		case "string":
			zapFields[i] = zap.String(f.Key, f.Value.(string))
		case "error":
			zapFields[i] = zap.Error(f.Value.(error))
		default:
			zapFields[i] = zap.Any(f.Key, f.Value)
		}
	}
	return zapFields
}

func (l *ZapLogger) Debug(msg string, fields ...Fields) {
	var zapFields []zap.Field
	if len(fields) == 1 {
		zapFields = l.FromFields(fields[0]).([]zap.Field)
	}
	l.inner.Debug(msg, zapFields...)
}

func (l *ZapLogger) Info(msg string, fields ...Fields) {
	var zapFields []zap.Field
	if len(fields) == 1 {
		zapFields = l.FromFields(fields[0]).([]zap.Field)
	}
	l.inner.Info(msg, zapFields...)
}

func (l *ZapLogger) Warn(msg string, fields ...Fields) {
	var zapFields []zap.Field
	if len(fields) == 1 {
		zapFields = l.FromFields(fields[0]).([]zap.Field)
	}
	l.inner.Warn(msg, zapFields...)
}

func (l *ZapLogger) Error(msg string, fields ...Fields) {
	var zapFields []zap.Field
	if len(fields) == 1 {
		zapFields = l.FromFields(fields[0]).([]zap.Field)
	}
	l.inner.Error(msg, zapFields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...Fields) {
	var zapFields []zap.Field
	if len(fields) == 1 {
		zapFields = l.FromFields(fields[0]).([]zap.Field)
	}
	l.inner.Fatal(msg, zapFields...)
}

func (l *ZapLogger) Panic(msg string, fields ...Fields) {
	var zapFields []zap.Field
	if len(fields) == 1 {
		zapFields = l.FromFields(fields[0]).([]zap.Field)
	}
	l.inner.Panic(msg, zapFields...)
}

// getZapLoggerLevel returns the minimum log level.
// If nothing is specified in the environment variable LOG_LEVEL,
// The level is DEBUG in development mode and WARN in others cases.
func getZapLoggerLevel(l string, env string) (level zapcore.Level) {
	switch l {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		if env == "development" {
			level = zapcore.DebugLevel
		} else {
			level = zapcore.WarnLevel
		}
	}
	return
}
