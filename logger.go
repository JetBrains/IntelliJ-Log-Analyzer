package main

import (
	"github.com/wailsapp/wails/v2/pkg/logger"
	"os"
	"time"
)

// AppLogger is a utility to log messages to a number of destinations
type AppLogger struct{}

// NewAppLogger creates a new Logger.
func NewAppLogger() logger.Logger {
	return &AppLogger{}
}

// Print works like Sprintf.
func (l *AppLogger) Print(message string) {
	println(message)
}

// Trace level logging. Works like Sprintf.
func (l *AppLogger) Trace(message string) {
	println(time.Now().Format("TRA | 15:04:05.000 | ") + message)
}

// Debug level logging. Works like Sprintf.
func (l *AppLogger) Debug(message string) {
	println(time.Now().Format("DEB | 15:04:05.000 | ") + message)
}

// Info level logging. Works like Sprintf.
func (l *AppLogger) Info(message string) {
	println(time.Now().Format("INF | 15:04:05.000 | ") + message)
}

// Warning level logging. Works like Sprintf.
func (l *AppLogger) Warning(message string) {
	println(time.Now().Format("WAR | 15:04:05.000 | ") + message)
}

// Error level logging. Works like Sprintf.
func (l *AppLogger) Error(message string) {
	println(time.Now().Format("ERR | 15:04:05.000 | ") + message)
}

// Fatal level logging. Works like Sprintf.
func (l *AppLogger) Fatal(message string) {
	println("FAT | " + message)
	os.Exit(1)
}
