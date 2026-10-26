package analyzer

import (
	"strings"
	"testing"
)

func TestProcessReaderStatisticsAndIncidents(t *testing.T) {
	input := strings.Join([]string{
		"2026-10-21T08:00:00Z INFO auth Login successful",
		"2026-10-21T08:00:01Z WARN auth Failed login attempt",
		"2026-10-21T08:00:02Z WARN auth Failed login attempt",
		"2026-10-21T08:00:03Z ERROR payment Payment processing failed",
		"2026-10-21T08:00:04Z ERROR payment Payment processing failed",
		"2026-10-21T08:00:05Z FATAL database Database connection lost",
		"INVALID LOG ENTRY",
	}, "\n")

	a := New()

	if err := ProcessReader(strings.NewReader(input), a); err != nil {
		t.Fatalf("ProcessReader returned an error: %v", err)
	}

	stats := a.Stats()

	if stats.TotalEntries != 6 {
		t.Errorf("expected 6 valid entries, got %d", stats.TotalEntries)
	}

	if stats.MalformedLines != 1 {
		t.Errorf("expected 1 malformed line, got %d", stats.MalformedLines)
	}

	if stats.CountBySeverity["INFO"] != 1 {
		t.Errorf("expected 1 INFO entry")
	}

	if stats.CountBySeverity["WARN"] != 2 {
		t.Errorf("expected 2 WARN entries")
	}

	if stats.CountBySeverity["ERROR"] != 2 {
		t.Errorf("expected 2 ERROR entries")
	}

	if stats.CountBySeverity["FATAL"] != 1 {
		t.Errorf("expected 1 FATAL entry")
	}

	if stats.CountByService["auth"] != 3 {
		t.Errorf("expected 3 auth entries")
	}

	expectedOccurrences := map[string]int{
		"Authentication Failure": 2,
		"Payment Failure":        2,
		"Database Outage":        1,
	}

	for _, incident := range stats.Incidents {
		expected, ok := expectedOccurrences[incident.Type]
		if !ok {
			t.Errorf("unexpected incident type %q", incident.Type)
			continue
		}

		if incident.Occurrences != expected {
			t.Errorf(
				"incident %q: expected %d occurrences, got %d",
				incident.Type,
				expected,
				incident.Occurrences,
			)
		}

		delete(expectedOccurrences, incident.Type)
	}

	for incidentType := range expectedOccurrences {
		t.Errorf("incident %q was not detected", incidentType)
	}
}

func TestProcessReaderOversizedLineIsSkipped(t *testing.T) {
	line := strings.Repeat("x", 1024*1024+1)
	input := line + "\n2026-10-21T08:00:00Z INFO auth Still processing"

	a := New()

	if err := ProcessReader(strings.NewReader(input), a); err != nil {
		t.Fatalf("ProcessReader returned an error: %v", err)
	}
	stats := a.Stats()
	if stats.MalformedLines != 1 {
		t.Errorf("expected 1 malformed line, got %d", stats.MalformedLines)
	}
	if stats.TotalEntries != 1 {
		t.Errorf("expected following valid line to be processed, got %d entries", stats.TotalEntries)
	}
}

func TestExpandedIncidentRulesAndOverlap(t *testing.T) {
	input := strings.Join([]string{
		"2026-10-21T08:00:00Z WARN auth Failed login attempt",
		"2026-10-21T08:00:01Z ERROR payment Payment database request failed",
		"2026-10-21T08:00:02Z ERROR auth NullPointerException in controller",
		"2026-10-21T08:00:03Z ERROR auth SQL syntax error in query",
		"2026-10-21T08:00:04Z WARN api High memory usage detected",
		"2026-10-21T08:00:05Z WARN api Database connection pool reaching threshold",
		"2026-10-21T08:00:06Z WARN edge Rate limit approaching",
		"2026-10-21T08:00:07Z WARN api Deprecated API endpoint accessed",
		"2026-10-21T08:00:08Z FATAL db Database connection lost",
	}, "\n")

	a := New()
	if err := ProcessReader(strings.NewReader(input), a); err != nil {
		t.Fatalf("ProcessReader returned an error: %v", err)
	}

	got := make(map[string]int)
	for _, incident := range a.Stats().Incidents {
		got[incident.Type+"/"+incident.Service] = incident.Occurrences
	}
	want := map[string]int{
		"Authentication Failure/auth":   1,
		"Database Connectivity/payment": 1,
		"Payment Failure/payment":       1,
		"Application Error/auth":        2,
		"Resource Pressure/api":         2,
		"Rate Limit Warning/edge":       1,
		"Deprecated API Usage/api":      1,
		"Database Outage/db":            1,
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d incident groups, got %d: %#v", len(want), len(got), got)
	}
	for key, expected := range want {
		if got[key] != expected {
			t.Errorf("%s: expected %d occurrences, got %d", key, expected, got[key])
		}
	}
}
