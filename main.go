package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"test_task_airports/core"
)

func main() {
	col, err := strconv.Atoi(os.Args[1])
	if err != nil || col < 1 {
		fmt.Println("Неверный номер колонки")
		return
	}

	col = col - 1

	parser := core.NewCsvParser()
	builder := core.NewCsvBuilder("./stuff/airports.csv", parser)

	index, err := builder.Build(col)
	if err != nil {
		fmt.Printf("Ошибка построения индекса - %s\n", err)
		return
	}

	searcher := core.NewCsvSearcher(index)

	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Введите текст для поиска")

		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Ошибка чтения ввода - %s\n", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == "!quit" {
			break
		}

		results := searcher.Search(input)
		for _, entry := range results {
			// Считываем соответствующую строку из файла по смещению
			lineBytes := make([]byte, entry.Length)
			_, err := index.File.Seek(entry.Offset, 0)
			if err != nil {
				continue
			}
			n, err := index.File.Read(lineBytes)
			if err != nil || n == 0 {
				continue
			}
			lineStr := strings.TrimRight(string(lineBytes), "\r\n")
			// Выводим найденное значение (без кавычек для строковых полей) и полную строку
			fmt.Printf("%s[%s]\n", entry.Key, lineStr)
		}
	}

	index.File.Close()
}
