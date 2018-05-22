package log

import (
    "os"

    "github.com/sirupsen/logrus"
)

func ConfigureDefault() {
    Configure(logrus.DebugLevel, &logrus.TextFormatter{})
}

func Configure(level logrus.Level, formatter logrus.Formatter) {
    logrus.SetOutput(os.Stdout)
    logrus.SetLevel(level)
    logrus.SetFormatter(formatter)
}

func WithField(key string, value interface{}) *logrus.Entry {
    return logrus.WithField(key, value)
}

type Fields logrus.Fields

func WithFields(fields Fields) *logrus.Entry {
    return logrus.WithFields(logrus.Fields(fields))
}

func WithError(err error) *logrus.Entry {
    return logrus.WithError(err)
}

func Info(msg interface{}) {
    logrus.Info(msg)
}

func Warn(msg interface{}) {
    logrus.Warn(msg)
}

func Error(msg interface{}) {
    logrus.Error(msg)
}

func Debug(msg interface{}) {
    logrus.Debug(msg)
}

func Fatal(msg interface{}) {
    logrus.Fatal(msg)
}