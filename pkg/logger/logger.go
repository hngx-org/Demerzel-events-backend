package logger

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var logger = logrus.New()

type Fields logrus.Fields

type Formatter struct {
	prefix string
}

func init() {
	logger.Level = logrus.InfoLevel
	logger.Formatter = &Formatter{}

	logger.SetReportCaller(true)
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString(" ")
	b.WriteString(entry.Time.Format(time.RFC3339))
	b.WriteString(" ")
	b.WriteString(f.prefix)
	b.WriteString(entry.Message)

	return b.Bytes(), nil
}

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

func Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatalf(format, args...)
	}
}
