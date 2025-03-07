package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogOutputsWithOneOutput(t *testing.T) {
	appName := "go-url-shortener"
	filePath := "/tmp"
	outputs := []string{"stdout"}

	gottenOutputs, err := getLoggerOutputs(outputs, "", "")
	assert.Equal(t, []string{"stdout"}, gottenOutputs, "with stdout")
	assert.Nil(t, err)

	outputs = []string{"file"}
	gottenOutputs, err = getLoggerOutputs(outputs, appName, filePath)
	assert.Equal(t, []string{"/tmp/go-url-shortener.log"}, gottenOutputs, "with file")
	assert.Nil(t, err)

	gottenOutputs, err = getLoggerOutputs(outputs, "", filePath)
	assert.Equal(t, []string(nil), gottenOutputs, "with file and empty app name")
	assert.NotNil(t, err)

	gottenOutputs, err = getLoggerOutputs(outputs, appName, "")
	assert.Equal(t, []string{"./go-url-shortener.log"}, gottenOutputs, "with file and empty file path")
	assert.Nil(t, err)
}

func TestLogOutputsWithMoreThanOneOutput(t *testing.T) {
	appName := "go-url-shortener"
	filePath := "/tmp"
	outputs := []string{"stdout"}

	gottenOutputs, err := getLoggerOutputs(outputs, "", "")
	assert.Equal(t, []string{"stdout"}, gottenOutputs, "with stdout")
	assert.Nil(t, err)

	outputs = []string{"stdout", "file"}
	gottenOutputs, err = getLoggerOutputs(outputs, appName, filePath)
	assert.Equal(t, []string{"/tmp/go-url-shortener.log", "stdout"}, gottenOutputs, "with stdout")
	assert.Nil(t, err)
}
