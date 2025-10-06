package logger

import (
	"log"
	"os"
	"strings"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel Level
	debugLogger  *log.Logger
	infoLogger   *log.Logger
	warnLogger   *log.Logger
	errorLogger  *log.Logger
)

func Init(level string) {
	debugLogger = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)
	infoLogger = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	warnLogger = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	switch strings.ToLower(level) {
	case "debug":
		currentLevel = DEBUG
	case "info":
		currentLevel = INFO
	case "warn":
		currentLevel = WARN
	case "error":
		currentLevel = ERROR
	default:
		currentLevel = INFO
	}
}

func Debug(format string, v ...interface{}) {
	if currentLevel <= DEBUG {
		debugLogger.Printf(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if currentLevel <= INFO {
		infoLogger.Printf(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if currentLevel <= WARN {
		warnLogger.Printf(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if currentLevel <= ERROR {
		errorLogger.Printf(format, v...)
	}
}
