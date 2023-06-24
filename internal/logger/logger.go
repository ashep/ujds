package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New() zerolog.Logger {
	l := log.Logger

	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice { // Terminal
		l = l.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return l
}
