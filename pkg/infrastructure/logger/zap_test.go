package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestZapLogLevel(t *testing.T) {
	cases := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"panic": zapcore.PanicLevel,
		"fatal": zapcore.FatalLevel,
		"":      zapcore.DebugLevel,
	}

	env := "development"
	for level, expected := range cases {
		assert.Equal(t, expected, getZapLoggerLevel(level, env))
	}

	env = "production"
	cases[""] = zapcore.WarnLevel
	for level, expected := range cases {
		assert.Equal(t, expected, getZapLoggerLevel(level, env))
	}
}

func TestZapFromFields(t *testing.T) {
	fields := Fields{
		NewField("code", "int", 500),
		NewField("message", "string", "Internal server error"),
		NewField("err", "error", errors.New("Internal server error")),
	}

	logger := &ZapLogger{}
	zapFields := logger.FromFields(fields).([]zapcore.Field)

	assert.Equal(t, 3, len(zapFields))
	assert.Equal(t, "code", zapFields[0].Key)
	assert.Equal(t, int64(500), zapFields[0].Integer)
	assert.Equal(t, "message", zapFields[1].Key)
	assert.Equal(t, "Internal server error", zapFields[1].String)
	assert.Equal(t, "error", zapFields[2].Key)
	assert.Equal(t, errors.New("Internal server error"), zapFields[2].Interface.(error))
}
