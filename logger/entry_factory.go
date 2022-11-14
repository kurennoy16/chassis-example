package logger

import (
	"net"
	"os"
)

type (
	EntryFunc func() Entry

	logEntryFactory struct {
		service string
	}

	LogEntryFactory interface {
		MakeEntry(level Level, code int, label string) EntryFunc
	}
)

func NewLogEntryFactory(service string) LogEntryFactory {
	return &logEntryFactory{service: service}
}

func (factory *logEntryFactory) MakeEntry(level Level, code int, label string) EntryFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "failed to determine"
	}

	return func() Entry {
		return newLogEntry(level, factory.service, code, label).
			WithField("hostname", hostname).
			WithField("address", ipAddress())
	}
}

func ipAddress() (address string) {
	address = ""
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			address += addr.String()
		}
	}

	return
}
