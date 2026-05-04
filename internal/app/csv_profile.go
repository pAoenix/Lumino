package app

import (
	"encoding/csv"
	"errors"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const previewLimit = 50

func ProfileCSV(path string) ([]string, int, []CSVRow, DataProfile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, nil, DataProfile{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if errors.Is(err, io.EOF) {
		return nil, 0, nil, DataProfile{}, nil
	}
	if err != nil {
		return nil, 0, nil, DataProfile{}, err
	}
	headers = normalizeHeaders(headers)

	stats := make(map[string]*numericAccumulator)
	preview := make([]CSVRow, 0, previewLimit)
	rows := 0

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return headers, rows, preview, DataProfile{}, err
		}
		rows++
		row := make(CSVRow, len(headers))
		for i, header := range headers {
			value := ""
			if i < len(record) {
				value = strings.TrimSpace(record[i])
			}
			row[header] = value
			if number, ok := parseNumber(value); ok {
				acc := stats[header]
				if acc == nil {
					acc = &numericAccumulator{min: math.Inf(1), max: math.Inf(-1)}
					stats[header] = acc
				}
				acc.add(number)
			}
		}
		if len(preview) < previewLimit {
			preview = append(preview, row)
		}
	}

	profile := DataProfile{Numeric: make([]NumericProfile, 0, len(stats))}
	for _, header := range headers {
		acc := stats[header]
		if acc == nil || acc.count == 0 {
			continue
		}
		profile.Numeric = append(profile.Numeric, NumericProfile{
			Column: header,
			Count:  acc.count,
			Min:    acc.min,
			Max:    acc.max,
			Avg:    acc.sum / float64(acc.count),
		})
	}

	return headers, rows, preview, profile, nil
}

type numericAccumulator struct {
	count int
	min   float64
	max   float64
	sum   float64
}

func (a *numericAccumulator) add(value float64) {
	a.count++
	a.sum += value
	a.min = math.Min(a.min, value)
	a.max = math.Max(a.max, value)
}

func normalizeHeaders(headers []string) []string {
	seen := make(map[string]int, len(headers))
	out := make([]string, len(headers))
	for i, header := range headers {
		header = strings.TrimSpace(header)
		if header == "" {
			header = "column_" + strconv.Itoa(i+1)
		}
		seen[header]++
		if seen[header] > 1 {
			header = header + "_" + strconv.Itoa(seen[header])
		}
		out[i] = header
	}
	return out
}

func parseNumber(value string) (float64, bool) {
	if value == "" {
		return 0, false
	}
	value = strings.ReplaceAll(value, ",", "")
	number, err := strconv.ParseFloat(value, 64)
	return number, err == nil
}
