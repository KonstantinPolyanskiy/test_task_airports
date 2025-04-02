package core

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"test_task_airports/common"
)

// Builder Строит индексы для файла, указанного в конструкторе builder'а-имплементации
type Builder interface {
	Build(columnNumber int) (common.Index, error)
}

type CsvBuilder struct {
	FilePath string
	// Построенный индекс
	Index common.Index

	Parser Parser
}

func NewCsvBuilder(filePath string, parser Parser) Builder {
	return CsvBuilder{
		FilePath: filePath,
		Parser:   parser,
	}
}

func (b CsvBuilder) Build(columnNumber int) (common.Index, error) {
	b.Index.Column = columnNumber

	f, err := os.Open(b.FilePath)
	if err != nil {
		return common.Index{}, err
	}
	b.Index.File = f

	reader := bufio.NewReader(f)
	var offset int64 = 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		lineLen := len(line)
		trimmedLine := strings.TrimRight(line, "\r\n")

		fields, parserErr := b.Parser.Parse(trimmedLine)
		if parserErr != nil || b.Index.Column < 0 || b.Index.Column >= len(fields) {
			fmt.Printf("Пропуск строки с offset=%d: parserErr=%v, fields=%v\n", offset, parserErr, fields)
			offset += int64(lineLen)
			if parserErr != nil {
				break
			}
			continue
		}

		normalString, isNum := b.toNormal(fields[b.Index.Column])

		//if len(b.Index.Entries) == 0 {
		//	b.Index.IsNumeric = isNum
		//}

		var numValue float64
		if isNum {
			numValue, _ = strconv.ParseFloat(normalString, 64)
		}

		entry := common.IndexEntry{
			Key:          normalString,
			NumericValue: numValue,
			IsNumeric:    isNum,
			Offset:       offset,
			Length:       lineLen,
		}

		b.Index.Entries = append(b.Index.Entries, entry)
		offset += int64(lineLen)
		if err == io.EOF {
			break
		}
	}

	sort.Slice(b.Index.Entries, func(i, j int) bool {
		return b.Index.Entries[i].Key < b.Index.Entries[j].Key
	})
	return b.Index, nil
}

// toNormal Вспомогательная функция, нормальзующая строку. Возвращает нормальизованную строку, и является ли строка числовым значением
func (b CsvBuilder) toNormal(original string) (string, bool) {
	var normalized string
	isNum := true

	if len(original) > 0 && original[0] == '"' && original[len(original)-1] == '"' {
		normalized = original[1 : len(original)-1]
	} else {
		normalized = original
	}

	normalized = strings.TrimSpace(normalized)
	normalized = strings.ToLower(normalized)

	// Если число
	if _, err := strconv.ParseFloat(normalized, 64); err != nil {
		isNum = false
	}

	return normalized, isNum
}
