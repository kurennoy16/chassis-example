package logger

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	Entry interface {
		Write(args ...interface{})
		Writeln(args ...interface{})
		Writef(format string, args ...interface{})
		WithField(key string, value interface{}) Entry
		WithFields(fields log.Fields) Entry
		Level() Level
		Entry() *log.Entry
	}

	entry struct {
		entry *log.Entry
		level Level
	}
)

func newLogEntry(level Level, service string, code int, label string) Entry {
	return &entry{
		entry: log.NewEntry(log.StandardLogger()).
			WithField("service", service).
			WithField("code", code).
			WithField("label", label),
		level: level,
	}
}

func (entry *entry) Level() Level {
	return entry.level
}

func (entry *entry) Entry() *log.Entry {
	return entry.entry
}

func (entry *entry) Write(args ...interface{}) {
	etry := addInstanceFields(entry, args)
	switch etry.Level() {
	case DEBUG:
		etry.Entry().Debug(args...)
	case ERROR:
		etry.Entry().Error(args...)
	case INFO:
		etry.Entry().Info(args...)
	}
}

func (entry *entry) Writeln(args ...interface{}) {
	etry := addInstanceFields(entry, args)
	switch etry.Level() {
	case DEBUG:
		etry.Entry().Debugln(args...)
	case ERROR:
		etry.Entry().Errorln(args...)
	case INFO:
		etry.Entry().Infoln(args...)
	}
}

func (entry *entry) Writef(format string, args ...interface{}) {
	etry := addInstanceFields(entry, args)
	switch etry.Level() {
	case DEBUG:
		etry.Entry().Debugf(format, args...)
	case ERROR:
		etry.Entry().Errorf(format, args...)
	case INFO:
		etry.Entry().Infof(format, args...)
	}
}

func (entry *entry) WithField(key string, value interface{}) Entry {
	entry.entry = entry.entry.WithField(key, value)
	return entry
}

func (entry *entry) WithFields(fields log.Fields) Entry {
	entry.entry = entry.entry.WithFields(fields)
	return entry
}

func addInstanceFields(entry *entry, args []interface{}) Entry {
	var etry Entry = entry
	if traceLevels != nil && ((*traceLevels & entry.level) != 0) {
		for index, arg := range args {
			if err, ok := arg.(error); ok {
				etry = etry.WithField(
					fmt.Sprintf("ErrorStack_%v", index),
					fmt.Sprintf("%+v", err))
			}
		}
	}

	return etry.WithField("epoch", time.Now().Unix())
}
