package commons

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
}

type AppLogger struct {
	logger *log.Logger
}

func NewLogger(prefix string) *AppLogger {
	return &AppLogger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags|log.Lshortfile),
	}
}

func (l *AppLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *AppLogger) Error(v ...interface{}) {
	l.logger.Printf("[ERROR] %s", fmt.Sprint(v...))
}

func (l *AppLogger) Errorf(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

func (l *AppLogger) Info(v ...interface{}) {
	l.logger.Printf("[INFO] %s", fmt.Sprint(v...))
}

func (l *AppLogger) Infof(format string, v ...interface{}) {
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *AppLogger) Debug(v ...interface{}) {
	l.logger.Printf("[DEBUG] %s", fmt.Sprint(v...))
}

func (l *AppLogger) Debugf(format string, v ...interface{}) {
	l.logger.Printf("[DEBUG] "+format, v...)
}

func (l *AppLogger) Warn(v ...interface{}) {
	l.logger.Printf("[WARN] %s", fmt.Sprint(v...))
}

func (l *AppLogger) Warnf(format string, v ...interface{}) {
	l.logger.Printf("[WARN] "+format, v...)
}
