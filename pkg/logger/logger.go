package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirawatc/simple-gin-crud/pkg/middleware"
	"github.com/sirupsen/logrus"
)

type customFormatter struct {
	format func(entry *logrus.Entry) ([]byte, error)
}

func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return f.format(entry)
}

func NewLogger(serviceName string) *logrus.Logger {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)

	logger.SetLevel(logrus.InfoLevel)

	logger.SetFormatter(&customFormatter{DefaultLogFormat(serviceName)})

	return logger
}

func DefaultLogFormat(serviceName string) func(entry *logrus.Entry) ([]byte, error) {
	return func(entry *logrus.Entry) ([]byte, error) {
		timestamp := entry.Time.In(time.FixedZone("GMT+7", 7*3600)).Format("2006-01-02T15:04:05Z07:00")
		logLevel := entry.Level.String()
		message := entry.Message
		requestId := entry.Data["requestId"]
		formattedMsg := fmt.Sprintf("[%s] [%s] [%s] : { requestId: %v, msg: %s }\n", serviceName, timestamp, logLevel, requestId, message)
		return []byte(formattedMsg), nil
	}
}

func InjectRequestIDWithLogger(ctx context.Context, logger *logrus.Logger) *logrus.Entry {
	return logger.WithField("requestId", middleware.GetRequestID(ctx))
}
