package logger

import (
	"fmt"
	"preza/internal/domain/logger"
	"strings"
	"time"
)

// Logger is the central logger for the application.
// Currently outputs to stdout. Replace the log() body to integrate
// with an external provider (Sentry, Datadog, etc.) without changing callers.
type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, fields ...logger.Field) {
	l.log("INFO", msg, fields)
}

func (l *Logger) Warn(msg string, fields ...logger.Field) {
	l.log("WARN", msg, fields)
}

func (l *Logger) Error(msg string, fields ...logger.Field) {
	l.log("ERROR", msg, fields)
}

func (l *Logger) log(level, msg string, fields []logger.Field) {
	var b strings.Builder
	fmt.Fprintf(&b, "%s [%s] %s", time.Now().UTC().Format(time.RFC3339), level, msg)
	for _, f := range fields {
		fmt.Fprintf(&b, " %s=%v", f.Key, f.Value)
	}
	fmt.Println(b.String())
}
