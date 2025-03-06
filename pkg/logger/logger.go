package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func InitLogger(defaultLevel string) error {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = defaultLevel
	}

	level, err := zerolog.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		return err
	}

	Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(level)

	return nil
}
