package main

import (
	"log"
	"os"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Purple  = "\033[35m"
	Cyan    = "\033[36m"
)

type Logger struct {
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	warnLogger    *log.Logger
	debugLogger   *log.Logger
	successLogger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLogger:    log.New(os.Stdout, Blue+"[INFO] "+Reset, log.Ldate|log.Ltime),
		errorLogger:   log.New(os.Stderr, Red+"[ERROR] "+Reset, log.Ldate|log.Ltime),
		warnLogger:    log.New(os.Stdout, Yellow+"[WARN] "+Reset, log.Ldate|log.Ltime),
		debugLogger:   log.New(os.Stdout, Cyan+"[DEBUG] "+Reset, log.Ldate|log.Ltime),
		successLogger: log.New(os.Stdout, Green+"[SUCCESS] "+Reset, log.Ldate|log.Ltime),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.warnLogger.Printf(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}

func (l *Logger) Success(format string, v ...interface{}) {
	l.successLogger.Printf(format, v...)
}
