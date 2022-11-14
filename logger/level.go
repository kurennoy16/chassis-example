package logger

import (
	"fmt"
	"strings"
)

const (
	DEBUG Level = 1 << iota
	INFO
	ERROR

	DefaultStackTraceLogLevel = DEBUG | ERROR
)

var traceLevels *Level

type Level byte

func ParseLevel(level string) (*Level, error) {
	switch strings.ToUpper(strings.Trim(level, "")) {
	case "DEBUG":
		return levelPtr(DEBUG), nil
	case "INFO":
		return levelPtr(INFO), nil
	case "ERROR":
		return levelPtr(ERROR), nil
	default:
		return nil, fmt.Errorf("unrecognised log level : %v", level)
	}
}

func SetStackTraceLevels(stackTraceLevel Level) {
	traceLevels = levelPtr(stackTraceLevel)
}

func levelPtr(level Level) *Level {
	return &level
}
