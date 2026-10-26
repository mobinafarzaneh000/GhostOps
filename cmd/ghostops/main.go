package main

import (
	"flag"
	"fmt"
	"os"

	"ghostops/internal/analyzer"
	"ghostops/internal/report"
)

func main() {
	outputPath := flag.String(
		"o",
		"",
		"Optional path for the JSON report",
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: ghostops [options] <log-file>\n\n")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	inputPath := flag.Arg(0)

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: input file path cannot be empty")
		os.Exit(2)
	}

	if *outputPath != "" && inputPath == *outputPath {
		fmt.Fprintln(
			os.Stderr,
			"Error: input and output paths must be different",
		)
		os.Exit(2)
	}

	if err := run(inputPath, *outputPath); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	a := analyzer.New()
	processErr := analyzer.ProcessReader(file, a)
	closeErr := file.Close()
	if processErr != nil {
		if closeErr != nil {
			return fmt.Errorf("process log file: %w (also failed to close input: %v)", processErr, closeErr)
		}
		return fmt.Errorf("process log file: %w", processErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close input file: %w", closeErr)
	}

	stats := a.Stats()

	if outputPath == "" {
		if err := report.WriteText(os.Stdout, stats); err != nil {
			return fmt.Errorf("write terminal report: %w", err)
		}

		return nil
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create JSON report: %w", err)
	}

	if err := report.WriteJSON(out, stats); err != nil {
		_ = out.Close()
		return fmt.Errorf("write JSON report: %w", err)
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("close JSON report: %w", err)
	}

	return nil
}
