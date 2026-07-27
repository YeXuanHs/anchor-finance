package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Entry for structured logging.
type Logger struct {
	*logrus.Entry
}

// New creates a configured Logger instance.
func New(level string) *Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)

	return &Logger{Entry: logrus.NewEntry(log)}
}

// WithField creates a child logger with a single field.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{Entry: l.Entry.WithField(key, value)}
}

// WithFields creates a child logger with multiple fields.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{Entry: l.Entry.WithFields(logrus.Fields(fields))}
}
