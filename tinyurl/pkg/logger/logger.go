package logger

import (
	"github.com/sirupsen/logrus"
)

var LoggerInstance Logger = logrus.New()

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func SetLogger(logger Logger) {
	LoggerInstance = logger
}

func Infof(format string, args ...interface{}) {
	LoggerInstance.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	LoggerInstance.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	LoggerInstance.Fatalf(format, args...)
}
