package analyzer

import (
	"sort"
	"strings"
	"time"

	"ghostops/internal/logparser"
)

type Incident struct {
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Service     string    `json:"service"`
	Occurrences int       `json:"occurrences"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
}

type incidentKey struct {
	Type    string
	Service string
}

// processIncidents checks an entry against the predefined rules.
func (a *Analyzer) processIncidents(entry logparser.LogEntry) {
	message := strings.ToLower(entry.Message)

	if entry.Level == "WARN" && strings.Contains(message, "failed login") {
		a.recordIncident(
			"Authentication Failure",
			"WARNING",
			entry,
		)
	}

	if entry.Level == "ERROR" && strings.Contains(message, "database") {
		a.recordIncident("Database Connectivity", "HIGH", entry)
	}

	if entry.Level == "ERROR" &&
		(strings.Contains(message, "nullpointerexception") ||
			strings.Contains(message, "sql syntax error")) {
		a.recordIncident("Application Error", "HIGH", entry)
	}

	if entry.Level == "ERROR" &&
		strings.Contains(message, "payment") {
		a.recordIncident(
			"Payment Failure",
			"HIGH",
			entry,
		)
	}

	if entry.Level == "WARN" &&
		(strings.Contains(message, "high memory usage") ||
			strings.Contains(message, "database connection pool reaching threshold")) {
		a.recordIncident("Resource Pressure", "MEDIUM", entry)
	}

	if entry.Level == "WARN" && strings.Contains(message, "rate limit approaching") {
		a.recordIncident("Rate Limit Warning", "MEDIUM", entry)
	}

	if entry.Level == "WARN" && strings.Contains(message, "deprecated api endpoint") {
		a.recordIncident("Deprecated API Usage", "LOW", entry)
	}

	if entry.Level == "FATAL" &&
		strings.Contains(message, "database") {
		a.recordIncident(
			"Database Outage",
			"CRITICAL",
			entry,
		)
	}
}

func (a *Analyzer) recordIncident(
	incidentType string,
	severity string,
	entry logparser.LogEntry,
) {
	key := incidentKey{
		Type:    incidentType,
		Service: entry.Service,
	}

	if index, ok := a.incidentIndexes[key]; ok {
		incident := &a.stats.Incidents[index]

		incident.Occurrences++

		if entry.Timestamp.Before(incident.FirstSeen) {
			incident.FirstSeen = entry.Timestamp
		}

		if entry.Timestamp.After(incident.LastSeen) {
			incident.LastSeen = entry.Timestamp
		}

		return
	}

	index := len(a.stats.Incidents)

	a.stats.Incidents = append(a.stats.Incidents, Incident{
		Type:        incidentType,
		Severity:    severity,
		Service:     entry.Service,
		Occurrences: 1,
		FirstSeen:   entry.Timestamp,
		LastSeen:    entry.Timestamp,
	})

	a.incidentIndexes[key] = index
}

func sortIncidents(incidents []Incident) {
	sort.Slice(incidents, func(i, j int) bool {
		a := incidents[i]
		b := incidents[j]

		if a.Type != b.Type {
			return a.Type < b.Type
		}

		return a.Service < b.Service
	})
}
