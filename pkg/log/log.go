package log

import "io"

const (
	DebugLevel string = "debug"
	InfoLevel  string = "info"
	WarnLevel  string = "warn"
	ErrorLevel string = "error"
)

var globalLogger Logger

func SetGlobalLogger(l Logger) {
	globalLogger = l
}

func GetGlobalLogger() Logger {
	return globalLogger
}

func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}
func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}
func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

type Logger interface {
	io.Writer
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
