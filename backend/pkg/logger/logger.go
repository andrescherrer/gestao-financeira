package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// InitLogger initializes the structured logger
func InitLogger(level string) {
	InitLoggerWithConfig(level, "", false)
}

// InitLoggerWithConfig initializes the structured logger with full configuration
func InitLoggerWithConfig(level, format string, isProduction bool) {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(parseLevel(level))

	// Configure output format
	output := os.Stdout
	useJSON := format == "json" || isProduction

	if useJSON {
		// In production or when format is json, use JSON format
		Logger = zerolog.New(output).
			With().
			Timestamp().
			Str("service", "gestao-financeira").
			Logger()
	} else {
		// In development, use console writer with colors
		Logger = zerolog.New(zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Str("service", "gestao-financeira").
			Logger()
	}

	// Set as global logger
	log.Logger = Logger
}

// parseLevel converts string level to zerolog.Level
func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

// GetLogger returns the configured logger instance
func GetLogger() zerolog.Logger {
	return Logger
}
