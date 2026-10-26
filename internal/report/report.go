package report

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"ghostops/internal/analyzer"
)

func WriteText(w io.Writer, stats analyzer.Stats) error {
	if _, err := fmt.Fprintln(w, "GhostOps Log Analysis"); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "====================="); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(
		w,
		"\nValid entries:   %d\nMalformed lines: %d\n",
		stats.TotalEntries,
		stats.MalformedLines,
	); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "\nSeverity Counts:"); err != nil {
		return err
	}

	severities := []string{"INFO", "WARN", "ERROR", "FATAL"}

	for _, severity := range severities {
		if _, err := fmt.Fprintf(
			w,
			"  %-6s %d\n",
			severity,
			stats.CountBySeverity[severity],
		); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w, "\nService Counts:"); err != nil {
		return err
	}

	services := make([]string, 0, len(stats.CountByService))

	for service := range stats.CountByService {
		services = append(services, service)
	}

	sort.Strings(services)

	for _, service := range services {
		if _, err := fmt.Fprintf(
			w,
			"  %-20s %d\n",
			service,
			stats.CountByService[service],
		); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w, "\nIncidents:"); err != nil {
		return err
	}

	if len(stats.Incidents) == 0 {
		_, err := fmt.Fprintln(w, "  No incidents detected.")
		return err
	}

	table := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	if _, err := fmt.Fprintln(
		table,
		"TYPE\tSEVERITY\tSERVICE\tCOUNT\tFIRST SEEN\tLAST SEEN",
	); err != nil {
		return err
	}

	for _, incident := range stats.Incidents {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%d\t%s\t%s\n",
			incident.Type,
			incident.Severity,
			incident.Service,
			incident.Occurrences,
			incident.FirstSeen.Format("2006-01-02T15:04:05Z07:00"),
			incident.LastSeen.Format("2006-01-02T15:04:05Z07:00"),
		); err != nil {
			return err
		}
	}

	return table.Flush()
}
