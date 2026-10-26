package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"ghostops/internal/analyzer"
)

func TestWriteJSON(t *testing.T) {
	stats := analyzer.Stats{
		TotalEntries:   2,
		MalformedLines: 1,
		CountBySeverity: map[string]int{
			"INFO": 2,
		},
		CountByService: map[string]int{
			"auth": 2,
		},
		Incidents: []analyzer.Incident{},
	}

	var output bytes.Buffer

	if err := WriteJSON(&output, stats); err != nil {
		t.Fatalf("WriteJSON returned an error: %v", err)
	}

	var decoded analyzer.Stats

	if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if decoded.TotalEntries != 2 {
		t.Errorf("expected 2 entries, got %d", decoded.TotalEntries)
	}

	if decoded.MalformedLines != 1 {
		t.Errorf(
			"expected 1 malformed line, got %d",
			decoded.MalformedLines,
		)
	}
}