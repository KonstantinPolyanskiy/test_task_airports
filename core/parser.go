package core

import "strings"

// Parser Парсит строку из CSV файла (с учетом ковычек внутри строк)
type Parser interface {
	Parse(line string) ([]string, error)
}

type CsvParser struct{}

func NewCsvParser() Parser {
	return CsvParser{}
}

func (p CsvParser) Parse(line string) ([]string, error) {
	var fields []string
	var currentFieldBuilder strings.Builder

	qoute := false

	for i := 0; i < len(line); i++ {
		char := line[i]
		if char == '"' {
			qoute = !qoute
			currentFieldBuilder.WriteByte(char)
		} else if char == ',' && !qoute {
			fields = append(fields, strings.TrimSpace(currentFieldBuilder.String()))
			currentFieldBuilder.Reset()
		} else {
			currentFieldBuilder.WriteByte(char)
		}
	}

	fields = append(fields, strings.TrimSpace(currentFieldBuilder.String()))

	return fields, nil
}
