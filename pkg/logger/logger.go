package logger

// DefaultInterface represents common interface
// for default logging function
type DefaultInterface interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
}

// FmtInterface represents common interface
// for formatting logging function
type FmtInterface interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// Interface represents application's logger contract
type Interface interface {
	DefaultInterface
	FmtInterface
}

// Logger is a struct that provides
// all functions for logging
type Logger struct {
	Interface
}

func New(newLogger Interface) *Logger {
	return &Logger{newLogger}
}
