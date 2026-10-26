package analyzer

import (
	"bufio"
	"fmt"
	"io"

	"ghostops/internal/logparser"
)

const maxLogLineSize = 1024 * 1024

func ProcessReader(reader io.Reader, a *Analyzer) error {
	buffered := bufio.NewReader(reader)
	line := make([]byte, 0, 64*1024)
	for {
		line = line[:0]
		tooLong := false

		for {
			part, err := buffered.ReadSlice('\n')
			if !tooLong {
				if len(line)+len(part) > maxLogLineSize {
					tooLong = true
					line = nil 
				} else {
					line = append(line, part...)
				}
			}

			if err == bufio.ErrBufferFull {
				continue
			}
			if err != nil && err != io.EOF {
				return fmt.Errorf("read log input: %w", err)
			}

			if err == io.EOF && len(part) == 0 && len(line) == 0 && !tooLong {
				return nil
			}

			if tooLong {
				a.RecordMalformed()
			} else {
				if len(line) > 0 && line[len(line)-1] == '\n' {
					line = line[:len(line)-1]
				}
				if len(line) > 0 && line[len(line)-1] == '\r' {
					line = line[:len(line)-1]
				}
				entry, parseErr := logparser.ParseLine(string(line))
				if parseErr != nil {
					a.RecordMalformed()
				} else {
					a.Process(entry)
				}
			}

			if err == io.EOF {
				return nil
			}
			break
		}
	}
}
