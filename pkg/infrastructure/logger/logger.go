package logger

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fabienbellanger/goutils"
)

// Fields is a slice of Field
type Fields = []Field

// Field is a struct that represents a log field.
type Field struct {
	Key   string
	Value any
	Type  string
}

// NewField creates a new field.
// Type must be one of the following: int, string, error.
func NewField(k string, t string, v any) Field {
	return Field{
		Key:   k,
		Type:  t,
		Value: v,
	}
}

// CustomLogger is the interface that a logger must implement.
type CustomLogger interface {
	FromFields(fields Fields) any
	Debug(msg string, fields ...Fields)
	Info(msg string, fields ...Fields)
	Warn(msg string, fields ...Fields)
	Error(msg string, fields ...Fields)
	Fatal(msg string, fields ...Fields)
	Panic(msg string, fields ...Fields)
}

// getLoggerOutputs returns an array with the log outputs.
// Outputs can be stdout and/or file.
func getLoggerOutputs(logOutputs []string, appName, filePath string) (outputs []string, err error) {
	if goutils.StringInSlice("file", logOutputs) {
		logPath := path.Clean(filePath)
		_, err := os.Stat(logPath)
		if err != nil {
			return nil, err
		}

		if appName == "" {
			return nil, errors.New("no APP_NAME variable defined")
		}

		outputs = append(outputs, fmt.Sprintf("%s/%s.log",
			logPath,
			appName))
	}
	if goutils.StringInSlice("stdout", logOutputs) {
		outputs = append(outputs, "stdout")
	}
	return
}
