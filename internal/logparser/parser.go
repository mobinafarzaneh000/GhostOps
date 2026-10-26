package logparser

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Service   string
	Message   string
}

var supportedLevels = map[string]bool{
	"INFO":  true,
	"WARN":  true,
	"ERROR": true,
	"FATAL": true,
}

func ParseLine(line string) (LogEntry, error) {
	timestampText, remainder, ok := takeField(line)
	if !ok {
		return LogEntry{}, fmt.Errorf("expected timestamp, level, service, and message")
	}

	level, remainder, ok := takeField(remainder)
	if !ok {
		return LogEntry{}, fmt.Errorf("expected timestamp, level, service, and message")
	}
	service, message, ok := takeField(remainder)
	if !ok {
		return LogEntry{}, fmt.Errorf("expected timestamp, level, service, and message")
	}

	timestamp, err := time.Parse(time.RFC3339, timestampText)
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	if !supportedLevels[level] {
		return LogEntry{}, fmt.Errorf("unsupported severity level %q", level)
	}

	if service == "" {
		return LogEntry{}, fmt.Errorf("service cannot be empty")
	}

	message = strings.TrimSpace(message)

	if message == "" {
		return LogEntry{}, fmt.Errorf("message cannot be empty")
	}

	return LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Service:   service,
		Message:   message,
	}, nil
}

func takeField(input string) (field, remainder string, ok bool) {
	input = strings.TrimLeftFunc(input, unicode.IsSpace)
	if input == "" {
		return "", "", false
	}
	end := strings.IndexFunc(input, unicode.IsSpace)
	if end == -1 {
		return input, "", true
	}
	return input[:end], input[end:], true
}
