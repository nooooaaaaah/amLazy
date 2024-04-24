package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from a .env file.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

// GetEnv retrieves the value of an environment variable.
func GetEnv(key string) string {
	return os.Getenv(key)
}

type CustomLogger struct {
	*log.Logger
	file *os.File // Keep a reference to the log file
}

func startLogger() *CustomLogger {
	logFile, err := os.OpenFile("amLazy.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	return &CustomLogger{Logger: logger, file: logFile}
}

func (l *CustomLogger) Close() {
	l.file.Close()
}

var startedLogger *CustomLogger // startedLogger holds a logger instance.

// InitLogger initializes the logger.
func InitLogger() {
	startedLogger = startLogger()
}

// GetLogger returns the current logger instance, initializing it if necessary.
func GetLogger() *CustomLogger {
	if startedLogger == nil {
		InitLogger()
	}
	return startedLogger
}

// LogInfo logs a message with an "INFO" prefix.
func (l *CustomLogger) LogInfo(msg string) {
	l.SetPrefix("INFO: ")
	l.Println(msg)
}

// LogError logs a message with an "ERROR" prefix.
func (l *CustomLogger) LogError(msg string) {
	l.SetPrefix("ERROR: ")
	l.Println(msg)
}

// LogErrorf logs a formatted error message.
func (l *CustomLogger) LogErrorf(format string, a ...interface{}) {
	l.SetPrefix("ERROR: ")
	l.Printf(format, a...)
}

// LogInfof logs a formatted info message.
func (l *CustomLogger) LogInfof(format string, a ...interface{}) {
	l.SetPrefix("INFO: ")
	l.Printf(format, a...)
}
