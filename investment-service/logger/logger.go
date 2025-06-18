package logger

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Sync() error
	With(fields ...Field) Logger
}

// Field is an alias to abstract the logger implementation
type Field = interface{}
