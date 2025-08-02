package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	serviceName := "test-service"
	logger := NewLogger(serviceName)

	assert.NotNil(t, logger)
	assert.Equal(t, logrus.InfoLevel, logger.GetLevel())
	assert.Equal(t, os.Stdout, logger.Out)
}

func TestLogger_WithFields(t *testing.T) {
	logger := NewLogger("test-service")

	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.WithFields(logrus.Fields{
		"user_id": "123",
		"action":  "login",
	}).Info("User logged in")

	output := buf.String()
	assert.Contains(t, output, "User logged in")
	assert.Contains(t, output, "test-service")
	assert.Contains(t, output, "requestId: <nil>")
	assert.Contains(t, output, "[info]")
}

func TestLogger_Levels(t *testing.T) {
	logger := NewLogger("test-service")
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	output := buf.String()
	assert.Contains(t, output, "Info message")
	assert.Contains(t, output, "Warning message")
	assert.Contains(t, output, "Error message")
	assert.Contains(t, output, "[info]")
	assert.Contains(t, output, "[warning]")
	assert.Contains(t, output, "[error]")
}

func TestLogger_WithRequestID(t *testing.T) {
	logger := NewLogger("test-service")
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.WithField("requestId", "req-123").Info("Request processed")

	output := buf.String()
	assert.Contains(t, output, "Request processed")
	assert.Contains(t, output, "requestId: req-123")
	assert.Contains(t, output, "test-service")
}

func TestLogger_WithoutRequestID(t *testing.T) {
	logger := NewLogger("author-service")
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.Info("Author created")

	output := buf.String()
	assert.Contains(t, output, "[author-service]")
	assert.Contains(t, output, "requestId: <nil>")
	assert.Contains(t, output, "msg: Author created")
}

func TestLogger_Formatting(t *testing.T) {
	logger := NewLogger("test-service")
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.Infof("%s created successfully", "test")

	output := buf.String()
	assert.Contains(t, output, "test created successfully")
	assert.Contains(t, output, "test-service")
	assert.Contains(t, output, "requestId: <nil>")
}
