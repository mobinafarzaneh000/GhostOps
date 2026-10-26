# GhostOps — Streaming Log Analysis CLI

GhostOps is a command-line application written in Go that analyzes application logs, calculates statistics, detects predefined incident patterns, and generates reports.

It uses only the Go standard library and processes logs incrementally to avoid loading the entire input file into memory.

## Features

- Streaming log processing.
- Validation of timestamps, severity levels, services, and messages.
- Counting of valid entries and malformed lines.
- Statistics grouped by severity and service.
- Detection and grouping of predefined incidents.
- Human-readable terminal reports.
- JSON report export.
- Error handling for invalid input and failed file operations.

## Requirements

- Go installed on your system.
- No external services or third-party dependencies.

## Build

From the project root:

```bash
go build -o ghostops ./cmd/ghostops
```

## Usage

Print a report to the terminal:

```bash
./ghostops testdata/sample.log
```

Export a report as JSON:

```bash
./ghostops -o report.json testdata/sample.log
```

Display command-line help:

```bash
./ghostops -h
```

## Input Format

Each log entry uses the following whitespace-separated format:

```text
TIMESTAMP LEVEL SERVICE MESSAGE
```

Example:

```text
2026-10-21T08:00:00Z INFO auth User login successful
2026-10-21T08:00:01Z WARN auth Failed login attempt
2026-10-21T08:00:03Z ERROR payment Payment processing failed
2026-10-21T08:00:05Z FATAL database Database connection lost
```

Timestamps must use RFC3339 format.

Supported severity levels:

- INFO
- WARN
- ERROR
- FATAL

Messages may contain spaces. The service name is a single whitespace-separated field.

Malformed entries are counted and skipped. Lines larger than 1 MiB are counted as malformed and discarded while processing continues. Other reader errors terminate processing with an error.

## Incident Detection

The initial implementation includes these predefined rules:

| Incident | Condition |
|---|---|
| Authentication Failure | `WARN` message contains `failed login` |
| Payment Failure | Severity is `ERROR` and the message contains `payment` |
| Database Connectivity | Severity is `ERROR` and the message contains `database` |
| Application Error | Severity is `ERROR` and the message contains `nullpointerexception` or `sql syntax error` |
| Resource Pressure | `WARN` message contains `high memory usage` or `database connection pool reaching threshold` |
| Rate Limit Warning | `WARN` message contains `rate limit approaching` |
| Deprecated API Usage | `WARN` message contains `deprecated api endpoint` |
| Database Outage | Severity is `FATAL` and the message contains `database` |

Message matching is case-insensitive.

Incidents are grouped by incident type and service across the entire input file. Their occurrence counts and timestamp ranges are updated as matching entries are processed.

Message matching is case-insensitive. An entry may match multiple incident rules (for example, an `ERROR` mentioning both payment and database), so incident occurrence totals may overlap and should not be interpreted as additional log entries.

## Architecture

- `cmd/ghostops`: command-line entry point and execution coordination.
- `internal/logparser`: log entry representation and validation.
- `internal/analyzer`: streaming processing, statistics, and incident detection.
- `internal/report`: terminal and JSON report generation.
- `testdata`: sample input logs.

The application separates parsing, analysis, and reporting so that each component can be tested independently.

## Memory Efficiency

The processor reads input incrementally using buffered I/O.

It retains aggregate statistics and incident groups rather than the complete input file or every parsed entry.

Memory usage therefore depends primarily on the size of the current line, the number of distinct services, and the number of incident groups.

## Error Handling

Malformed lines are skipped and counted.

Fatal execution errors, such as input read failures or report output failures, are returned to the CLI and result in a nonzero exit status.

## Testing

Run all tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Run the race detector:

```bash
go test -race ./...
```

## Dependencies

GhostOps uses only the Go standard library.
