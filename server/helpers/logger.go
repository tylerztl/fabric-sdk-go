package helpers

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
}

func GetLogger() *logrus.Logger {
	return logger
}
