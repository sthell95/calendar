package logger

import (
	"github.com/sirupsen/logrus"
)

type Level string

const Error Level = "error"
const Warning Level = "warning"
const Info Level = "info"

var _ Logger = (*Log)(nil)

type Logger interface {
	Write(level Level, message string)
}

type Log struct{}

func (_ *Log) Write(level Level, message string) {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "02-01-2006 15:04-05"
	logrus.SetFormatter(formatter)
	formatter.FullTimestamp = true

	switch level {
	case Error:
		logrus.Fatal(message)
	case Warning:
		logrus.Warning(message)
	case Info:
		logrus.Info(message)
	}
}

func NewLogger() Logger {
	return &Log{}
}
