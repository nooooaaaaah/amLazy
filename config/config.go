package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "amLazy", "config.env")
	if err := godotenv.Load(configPath); err != nil {
		log.Printf("No .env file found at %s", configPath)
	}
}

func GetEnv(key string) string {
	ev := os.Getenv(key)
	if ev == "" {
		log.Fatalf("The value for %s was not found in environment", key)
	}
	return ev
}

func getUserShell() string {
	return GetEnv("USERS_SHELL")
}

func getUserOS() string {
	return GetEnv("USERS_OS")
}

func GetEnvInstructions() string {
	return fmt.Sprintf("My shell is %s, and my os is %s. Only return one option, just return the command", getUserShell(), getUserOS())
}

type CustomLogger struct {
	*log.Logger
	file    *os.File // Keep a reference to the log file
	enabled bool     // Flag to enable or disable logging
}

func startLogger() *CustomLogger {
	loggingEnabled := os.Getenv("LOGGING_ENABLED") // Check the environment variable
	enableLogging := loggingEnabled == "true"

	logFile, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".config", "amLazy", "amLazy.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	return &CustomLogger{Logger: logger, file: logFile, enabled: enableLogging}
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

// LogInfo logs a message with an "INFO" prefix if logging is enabled.
func (l *CustomLogger) LogInfo(msg string) {
	if l.enabled {
		l.SetPrefix("INFO: ")
		l.Println(msg)
	}
}

// LogError logs a message with an "ERROR" prefix if logging is enabled.
func (l *CustomLogger) LogError(msg string) {
	if l.enabled {
		l.SetPrefix("ERROR: ")
		l.Println(msg)
	}
}

// LogErrorf logs a formatted error message if logging is enabled.
func (l *CustomLogger) LogErrorf(format string, a ...interface{}) {
	if l.enabled {
		l.SetPrefix("ERROR: ")
		l.Printf(format, a...)
	}
}

// LogInfof logs a formatted info message if logging is enabled.
func (l *CustomLogger) LogInfof(format string, a ...interface{}) {
	if l.enabled {
		l.SetPrefix("INFO: ")
		l.Printf(format, a...)
	}
}
