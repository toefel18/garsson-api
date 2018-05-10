package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

func ConfigureDefault() {
	Configure(logrus.DebugLevel, &logrus.TextFormatter{})
}

func Configure(level logrus.Level, formatter logrus.Formatter) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(level)
	logrus.SetFormatter(formatter)
}

