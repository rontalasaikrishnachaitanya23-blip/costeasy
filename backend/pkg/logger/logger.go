// backend/pkg/logger/logger.go
package logger

import (
	"log"
	"os"
)

// Logger wraps standard logger
type Logger struct {
	*log.Logger
}

// New creates a new logger
func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// Info logs info message
func (l *Logger) Info(msg string) {
	l.Printf("[INFO] %s\n", msg)
}

// Error logs error message
func (l *Logger) Error(msg string) {
	l.Printf("[ERROR] %s\n", msg)
}

// Debug logs debug message
func (l *Logger) Debug(msg string) {
	l.Printf("[DEBUG] %s\n", msg)
}

// Warn logs warning message
func (l *Logger) Warn(msg string) {
	l.Printf("[WARN] %s\n", msg)
}
