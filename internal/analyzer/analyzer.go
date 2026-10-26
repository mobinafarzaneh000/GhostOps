package analyzer

import "ghostops/internal/logparser"

type Stats struct {
	TotalEntries    int            `json:"total_entries"`
	MalformedLines  int            `json:"malformed_lines"`
	CountBySeverity map[string]int `json:"count_by_severity"`
	CountByService  map[string]int `json:"count_by_service"`
	Incidents       []Incident     `json:"incidents"`
}

type Analyzer struct {
	stats          Stats
	incidentIndexes map[incidentKey]int
}

func New() *Analyzer {
	return &Analyzer{
		stats: Stats{
			CountBySeverity: make(map[string]int),
			CountByService:  make(map[string]int),
			Incidents:       []Incident{},
		},
		incidentIndexes: make(map[incidentKey]int),
	}
}

func (a *Analyzer) Process(entry logparser.LogEntry) {
	a.stats.TotalEntries++

	a.stats.CountBySeverity[entry.Level]++
	a.stats.CountByService[entry.Service]++

	a.processIncidents(entry)
}

func (a *Analyzer) RecordMalformed() {
	a.stats.MalformedLines++
}

func (a *Analyzer) Stats() Stats {
	stats := a.stats

	stats.CountBySeverity = make(map[string]int, len(a.stats.CountBySeverity))
	for key, value := range a.stats.CountBySeverity {
		stats.CountBySeverity[key] = value
	}

	stats.CountByService = make(map[string]int, len(a.stats.CountByService))
	for key, value := range a.stats.CountByService {
		stats.CountByService[key] = value
	}

	stats.Incidents = append([]Incident(nil), a.stats.Incidents...)
	sortIncidents(stats.Incidents)

	return stats
}
