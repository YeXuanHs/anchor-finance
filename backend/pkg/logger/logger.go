package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Config holds logger configuration (read from database).
type Config struct {
	Level  string
	Format string
	Output string
	File   string
}

// Logger wraps logrus.Entry for structured logging.
type Logger struct {
	*logrus.Entry
}

var defaultLogger *Logger

// Init initializes the global logger from database config.
func Init(cfg Config) {
	log := logrus.New()

	// Output
	if cfg.Output == "file" && cfg.File != "" {
		f, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(f)
		}
	} else {
		log.SetOutput(os.Stdout)
	}

	// Format
	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Level
	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)

	defaultLogger = &Logger{Entry: logrus.NewEntry(log)}
}

// New creates a configured Logger instance (for backward compatibility).
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

// Default returns the global logger instance.
func Default() *Logger {
	if defaultLogger == nil {
		return New("info")
	}
	return defaultLogger
}

// Infof logs at info level.
func Infof(format string, args ...interface{}) {
	Default().Infof(format, args...)
}

// Warnf logs at warn level.
func Warnf(format string, args ...interface{}) {
	Default().Warnf(format, args...)
}

// Errorf logs at error level.
func Errorf(format string, args ...interface{}) {
	Default().Errorf(format, args...)
}

// Info logs at info level.
func Info(args ...interface{}) {
	Default().Info(args...)
}

// WithField creates a child logger with a single field.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{Entry: l.Entry.WithField(key, value)}
}

// WithFields creates a child logger with multiple fields.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{Entry: l.Entry.WithFields(logrus.Fields(fields))}
}
