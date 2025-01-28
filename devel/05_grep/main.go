package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type parametres struct {
	after      int
	before     int
	toCount    bool
	ignoreCase bool
	invert     bool
	fixed      bool
	withLine   bool

	pattern  string
	filename string
}

var params parametres

func main() {
	params = parseArgsIntoParams()
	data, err := readData(params.filename)
	if err != nil {
		log.Fatalln(err)
	}

	doGrep(data, params, os.Stdout)
}

func parseArgsIntoParams() parametres {
	linesAfter := flag.Int("A", 0, "Print N lines after")
	linesBefore := flag.Int("B", 0, "Print N lines before")
	linesNear := flag.Int("C", 0, "Print N lines before and after")
	toCount := flag.Bool("c", false, "Only count lines")
	ignore := flag.Bool("i", false, "Ignore case ")
	invert := flag.Bool("v", false, "Invert results")
	fixed := flag.Bool("F", false, "Fixed pattern")
	withLine := flag.Bool("n", false, "Print lines with number")

	flag.Parse()

	params := parametres{
		after:      *linesAfter,
		before:     *linesBefore,
		toCount:    *toCount,
		ignoreCase: *ignore,
		invert:     *invert,
		fixed:      *fixed,

		withLine: *withLine,
	}

	if *linesNear != 0 {
		params.after = *linesNear
		params.before = *linesNear
	}

	if len(flag.Args()) != 2 {
		return params
	}
	params.pattern = flag.Args()[0]
	params.filename = flag.Args()[1]

	return params
}

func doGrep(data []string, params parametres, out io.Writer) {
	// Находим нужные индексы строк для вывода
	indices := findIndices(data, params)

	// добавляем строки сверху/снизу по ключам -A -B -C
	indices = addAreaIndices(indices, params, len(data))

	// Выводим результат
	printResult(data, indices, out, params.toCount)
}

func printResult(data []string, indices []int, out io.Writer, isCountNeed bool) {
	if isCountNeed {
		fmt.Fprintf(out, "%d\n", len(indices))
		return
	}

	for i, index := range indices {
		if i != len(indices)-1 {
			fmt.Fprintf(out, "%s\n", data[index])
		} else {
			fmt.Fprintf(out, "%s", data[index])
		}
	}
}

// проверяем за счёт текста
func checkByRaw(checked string) bool {
	var ans bool
	if params.ignoreCase {
		ans = strings.Contains(strings.ToLower(checked), strings.ToLower(params.pattern))
	} else {
		ans = strings.Contains(checked, params.pattern)
	}
	return ans
}

// проверяем на счёт регекса
func checkByRegex(checked string) bool {
	var ans bool
	pattern := params.pattern
	if params.ignoreCase {
		pattern = "(?i)" + pattern
	}
	ans, _ = regexp.MatchString(pattern, checked)
	return ans
}

func findIndices(data []string, params parametres) []int {
	var indices []int

	inverter := func(v bool) bool {
		if params.invert {
			return !v
		}
		return v
	}

	var checker func(string) bool
	if params.fixed {
		checker = checkByRaw
	} else {
		checker = checkByRegex
	}

	for i := 0; i < len(data); i++ {
		if inverter(checker(data[i])) {
			indices = append(indices, i)
		}
	}

	return indices
}

func addAreaIndices(inds []int, params parametres, maxIndex int) []int {
	if params.after == 0 && params.before == 0 {
		return inds
	}

	unique := make(map[int]struct{})
	for _, index := range inds {
		unique[index] = struct{}{}

		for i := 1; i <= params.before; i++ {
			newIndex := index - i
			if newIndex >= 0 {
				unique[newIndex] = struct{}{}
			}
		}

		for i := 1; i <= params.after; i++ {
			newIndex := index + i
			if newIndex < maxIndex {
				unique[newIndex] = struct{}{}
			}
		}
	}

	newIndices := make([]int, 0, len(unique))
	for key := range unique {
		newIndices = append(newIndices, key)
	}
	sort.Ints(newIndices)

	return newIndices
}

func readData(filename string) ([]string, error) {
	var input io.Reader
	if filename == "" {
		input = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return []string{}, err
		}
		defer file.Close()
	}

	scanner := bufio.NewScanner(input)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}
	return lines, nil
}
