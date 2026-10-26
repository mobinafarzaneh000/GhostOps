package report

import (
	"encoding/json"
	"io"

	"ghostops/internal/analyzer"
)

func WriteJSON(w io.Writer, stats analyzer.Stats) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(stats); err != nil {
		return err
	}

	return nil
}
