package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger provides structured logging capabilities
type Logger struct {
	prefix string
}

// New creates a new Logger instance
func New(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

// Default returns a logger with no prefix
func Default() *Logger {
	return &Logger{prefix: ""}
}

// Info logs informational messages
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.log("INFO", msg, fields)
}

// Error logs error messages
func (l *Logger) Error(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	l.log("ERROR", msg, fields)
}

// Warn logs warning messages
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.log("WARN", msg, fields)
}

// Debug logs debug messages
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	if os.Getenv("DEBUG") == "true" {
		l.log("DEBUG", msg, fields)
	}
}

// log is the internal logging method
func (l *Logger) log(level, msg string, fields map[string]interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	prefix := l.prefix
	if prefix != "" {
		prefix = fmt.Sprintf("[%s] ", prefix)
	}
	
	fieldStr := ""
	if len(fields) > 0 {
		fieldStr = fmt.Sprintf(" %+v", fields)
	}
	
	log.Printf("%s [%s] %s%s%s\n", timestamp, level, prefix, msg, fieldStr)
}

// WithPrefix returns a new logger with the specified prefix
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{prefix: prefix}
}
