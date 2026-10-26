package logparser

import (
	"testing"
	"time"
)

func TestParseLineValid(t *testing.T) {
	line := "2026-10-21T08:00:00Z INFO auth User login successful"

	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("ParseLine returned an unexpected error: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("expected INFO, got %q", entry.Level)
	}

	if entry.Service != "auth" {
		t.Errorf("expected auth, got %q", entry.Service)
	}

	if entry.Message != "User login successful" {
		t.Errorf("unexpected message: %q", entry.Message)
	}

	expectedTime := time.Date(
		2026, 10, 21, 8, 0, 0, 0, time.UTC,
	)

	if !entry.Timestamp.Equal(expectedTime) {
		t.Errorf(
			"expected timestamp %v, got %v",
			expectedTime,
			entry.Timestamp,
		)
	}
}

func TestParseLinePreservesMessageSpacing(t *testing.T) {
	line := "2026-10-21T08:00:00Z INFO auth  request   timed out  "
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("ParseLine returned an unexpected error: %v", err)
	}
	if entry.Message != "request   timed out" {
		t.Errorf("expected message spacing to be preserved, got %q", entry.Message)
	}
}

func TestParseLineInvalid(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{
			name: "too few fields",
			line: "2026-10-21T08:00:00Z INFO",
		},
		{
			name: "invalid timestamp",
			line: "not-a-timestamp INFO auth Login successful",
		},
		{
			name: "unsupported severity",
			line: "2026-10-21T08:00:00Z DEBUG auth Debug message",
		},
		{
			name: "missing message",
			line: "2026-10-21T08:00:00Z INFO auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseLine(tt.line)

			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	}
}
