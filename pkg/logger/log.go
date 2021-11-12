package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Level string

const (
	Error   Level = "error"
	Warning Level = "warning"
	Info    Level = "info"
)

var _ Logger = (*Log)(nil)

type Logger interface {
	Write(level Level, message string, code string)
}

type Log struct{}

func (*Log) Write(level Level, message string, code string) {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "02-01-2006 15:04:05"
	logrus.SetFormatter(formatter)
	formatter.FullTimestamp = true

	message = fmt.Sprintf("[%s] %s", code, message)
	switch level {
	case Error:
		logrus.Error(message)
	case Warning:
		logrus.Warning(message)
	case Info:
		logrus.Info(message)
	}
}

func NewLogger() Logger {
	return &Log{}
}
