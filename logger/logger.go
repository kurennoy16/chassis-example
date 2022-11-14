package logger

import (
	"io"

	log "github.com/sirupsen/logrus"
)

func Init(level Level, out io.Writer, pretty bool) {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: pretty})

	log.SetOutput(out)

	setLogLevel(level)

	traceLevels = levelPtr(DefaultStackTraceLogLevel)
}

func setLogLevel(level Level) {
	switch level {
	case DEBUG:
		log.SetLevel(log.DebugLevel)
	case INFO:
		log.SetLevel(log.InfoLevel)
	case ERROR:
		log.SetLevel(log.ErrorLevel)
	}
}
