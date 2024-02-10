package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(debug bool) *logrus.Logger {
	var loglevel logrus.Level
	if debug {
		loglevel = logrus.DebugLevel
	} else {
		loglevel = logrus.InfoLevel
	}
	formatter := logrus.TextFormatter{
		ForceColors: true,
	}
	return &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &formatter,
		Hooks:     make(logrus.LevelHooks),
		Level:     loglevel,
	}
}
