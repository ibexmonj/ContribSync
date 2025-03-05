package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

// InitLogger initializes a structured logger for stdout
func InitLogger(logLevel string) error {
	// Configure Zerolog to write logs in JSON format to stdout
	zerolog.TimeFieldFormat = time.RFC3339
	Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "csync").
		Logger()

	// Set log level based on the provided configuration
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)

	return nil
}
