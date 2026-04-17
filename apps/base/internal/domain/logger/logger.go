package logger

type Logger interface {
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}

func F(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// ErrorKind classifies errors so external providers (Sentry, etc.) can filter
// noise: validation errors are expected, system errors warrant alerts.
type ErrorKind string

const (
	KindNotFound   ErrorKind = "not_found"
	KindValidation ErrorKind = "validation"
	KindSystem     ErrorKind = "system"
)
