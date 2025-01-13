package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type parameters struct {
	column    int
	byNumeric bool
	isReverse bool
	isUnique  bool
	filenames []string
}

var params parameters

func parseArgsIntoParams() parameters {
	colNum := flag.Int("k", 1, "Number of column to sort")
	byNum := flag.Bool("n", false, "Numeric compare")
	isUn := flag.Bool("u", false, "Unique strings")
	isRev := flag.Bool("r", false, "Reverse order")
	flag.Parse()

	params := parameters{
		column:    *colNum,
		byNumeric: *byNum,
		isReverse: *isRev,
		isUnique:  *isUn,
		filenames: flag.Args(),
	}

	return params
}

// readDataFromFiles поддерживает чтение из нескольких файлов, указанных в параметрах, по аналогии с консольной утилитой sort
func readDataFromFiles(filenames []string) ([]string, error) {
	readDataFromFile := func(name string) ([]string, error) {
		file, err := os.Open(name)
		if err != nil {
			return []string{}, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lines := []string{}
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return []string{}, err
		}
		return lines, nil
	}

	var data []string
	for _, name := range filenames {
		newLines, err := readDataFromFile(name)
		if err != nil {
			return []string{}, err
		}
		data = append(data, newLines...)
	}

	return data, nil
}

func compareAsNumbers(lhs, rhs string) bool {
	lnum, lerr := strconv.Atoi(lhs)
	rnum, rerr := strconv.Atoi(rhs)

	//сравниваем как обычные строки
	if lerr != nil && rerr != nil {
		return lhs < rhs
	}
	// строки без числовых значений в указанной колонке выводятся в самом конце
	if lerr != nil || rerr != nil {
		return lerr == nil
	}
	return lnum < rnum
}

func compareAsStrings(lhs, rhs string) bool {
	return lhs < rhs
}

func doSort(data []string, params parameters) []string {
	if params.isUnique {
		data = makeStringsUnique(data)
	}

	var valueComparator func(string, string) bool
	if params.byNumeric {
		valueComparator = compareAsNumbers
	} else {
		valueComparator = compareAsStrings
	}

	// логика сравнения:
	// осуществляется проверка существования указанной колонки, сначала выводятся строки, у которых ее не существует,
	// сравниваются по дефолтной первой. Функция сравнения определяется в зависимости от наличия ключа -n
	compareLogic := func(i, j int) bool {
		lhs := strings.Split(data[i], " ")
		rhs := strings.Split(data[j], " ")
		if len(lhs) == 0 {
			return true
		}
		if len(rhs) == 0 {
			return false
		}

		if len(lhs) < params.column && len(rhs) >= params.column {
			return true
		}
		if len(lhs) >= params.column && len(rhs) < params.column {
			return false
		}

		if len(lhs) < params.column && len(rhs) < params.column {
			return valueComparator(lhs[0], rhs[0])
		}
		if len(lhs) >= params.column && len(rhs) >= params.column {
			return valueComparator(lhs[params.column-1], rhs[params.column-1])
		}
		// Ошибка логики программы. Используется именно паника, так как это скрипт, здесь можно просто выйти из программы
		panic("DEBUG: code should not run here")
	}

	if params.isReverse {
		sort.Slice(data, func(i, j int) bool {
			return !compareLogic(i, j)
		})
	} else {
		sort.Slice(data, compareLogic)
	}

	return data
}

func makeStringsUnique(strs []string) []string {
	keys := make(map[string]struct{}, len(strs))
	newData := []string{}

	for i := 0; i < len(strs); i++ {
		if _, ok := keys[strs[i]]; !ok {
			keys[strs[i]] = struct{}{}
			newData = append(newData, strs[i])
		}
	}
	return newData
}

func main() {
	// пример ввода параметров: go run main.go -k 2 -n file.txt numeric.txt
	// в данном примере произойдет сортировка данных из двух файлов с параметрами: числовая сортировка данных второй колонки
	params = parseArgsIntoParams()
	data, err := readDataFromFiles(params.filenames)
	if err != nil {
		log.Fatalln(err)
	}

	data = doSort(data, params)
	for _, str := range data {
		fmt.Println(str)
	}
}
