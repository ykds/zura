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
func Debugw(msg string, kvs ...interface{}) {
	globalLogger.Debugw(msg, kvs...)
}
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}
func Infow(msg string, kvs ...interface{}) {
	globalLogger.Infow(msg, kvs...)
}
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}
func Warnw(msg string, kvs ...interface{}) {
	globalLogger.Warnw(msg, kvs...)
}
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}
func Errorw(msg string, kvs ...interface{}) {
	globalLogger.Errorw(msg, kvs...)
}
func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}
func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}
func Panicw(msg string, kvs ...interface{}) {
	globalLogger.Panicw(msg, kvs...)
}
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}
func Fatalw(msg string, kvs ...interface{}) {
	globalLogger.Fatalw(msg, kvs...)
}

type Logger interface {
	io.Writer
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugw(msg string, kvs ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infow(msg string, kvs ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnw(msg string, kvs ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorw(msg string, kvs ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicw(msg string, kvs ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalw(msg string, kvs ...interface{})
}
