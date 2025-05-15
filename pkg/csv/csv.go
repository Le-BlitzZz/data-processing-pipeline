package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type HandlerFunc func(headers, record []string)

func ForEachRow(filePath string, fn HandlerFunc) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read headers %s: %w", filePath, err)
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			continue
		}

		fn(headers, record)
	}
}
