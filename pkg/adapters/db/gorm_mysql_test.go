package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestGetGormLogLevel(t *testing.T) {
	type input struct {
		level string
		env   string
	}

	tests := []struct {
		name   string
		args   input
		wanted logger.LogLevel
	}{
		{
			name:   "Silent level",
			args:   input{level: "silent", env: "development"},
			wanted: logger.Silent,
		},
		{
			name:   "Info level",
			args:   input{level: "info", env: "development"},
			wanted: logger.Info,
		},
		{
			name:   "Warn level",
			args:   input{level: "warn", env: "development"},
			wanted: logger.Warn,
		},
		{
			name:   "Error level",
			args:   input{level: "error", env: "development"},
			wanted: logger.Error,
		},
		{
			name:   "Unknown level in development environment",
			args:   input{level: "fatal", env: "development"},
			wanted: logger.Warn,
		},
		{
			name:   "Unknown level in production environment",
			args:   input{level: "panic", env: "production"},
			wanted: logger.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := getGormLogLevel(tt.args.level, tt.args.env)

			assert.Equal(t, level, tt.wanted)
		})
	}
}
