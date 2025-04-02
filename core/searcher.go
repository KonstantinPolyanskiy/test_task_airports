package core

import (
	"fmt"
	"sort"
	"strings"
	"test_task_airports/common"
	"time"
)

// Searcher Ищет в common.IndexEntry подходящие записи по префиксу без учета регистра
type Searcher interface {
	Search(prefix string) []common.IndexEntry
}

// CsvSearcher Сущность для поиска в индексированных данных
type CsvSearcher struct {
	index common.Index
}

func NewCsvSearcher(index common.Index) Searcher {
	return CsvSearcher{index: index}
}

func (s CsvSearcher) Search(prefix string) []common.IndexEntry {
	start := time.Now()

	prefix = strings.ToLower(prefix)

	val := sort.Search(len(s.index.Entries), func(i int) bool {
		return s.index.Entries[i].Key >= prefix
	})

	var result []common.IndexEntry

	for i := val; i < len(s.index.Entries); i++ {
		if strings.HasPrefix(strings.ToLower(s.index.Entries[i].Key), prefix) {
			result = append(result, s.index.Entries[i])
		} else {
			break
		}
	}

	if s.index.IsNumeric {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Key < result[j].Key
		})
	}

	elapsed := time.Since(start)
	fmt.Printf("Найдено совпадений: %d, время поиска: %d\n", len(result), elapsed.Microseconds())

	return result
}
