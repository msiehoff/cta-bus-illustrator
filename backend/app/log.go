package app

import (
	"log"
	"os"
	"strings"
)

var debugLogging bool

// InitLogLevel reads LOG_LEVEL from the environment.
// Supported values: "info" (default), "debug".
func InitLogLevel() {
	level := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	debugLogging = level == "debug"
	if debugLogging {
		log.Println("log level: debug")
	}
}

func Debugf(format string, args ...any) {
	if debugLogging {
		log.Printf("DEBUG "+format, args...)
	}
}

func DebugEnabled() bool {
	return debugLogging
}
