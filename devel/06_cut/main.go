package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type parametres struct {
	fields        []int
	delim         string
	withDelimOnly bool

	filename string
}

func main() {
	// Парсим входные аргументы
	params, err := parseArgsIntoParams()
	if err != nil {
		fmt.Println("Invalid arg!", err.Error())
		return
	}

	data, err := readData(params.filename)
	if err != nil {
		fmt.Println("input error: ", err.Error())
		return
	}

	doCut(data, params, os.Stdout)
}

func parseArgsIntoParams() (parametres, error) {
	fields := flag.String("f", "", "Choose fileds")
	delim := flag.String("d", "", "Choose delim")
	isSep := flag.Bool("s", false, "Separated parametr")

	flag.Parse()

	// Для флага -f нужно извлечь допонительные данные
	// Сразу же проверяем уникальность параметров от пользователя
	uniqueColumns := make(map[int]struct{})
	for _, column := range strings.Split(*fields, " ") {
		num, err := strconv.Atoi(column)
		if err != nil {
			return parametres{}, fmt.Errorf("column -f error")
		}
		if num < 1 {
			return parametres{}, fmt.Errorf("Field must be at least 1.")
		}
		uniqueColumns[num] = struct{}{}
	}
	columns := make([]int, 0, len(uniqueColumns))
	for unColumn := range uniqueColumns {
		columns = append(columns, unColumn-1)
	}
	sort.Ints(columns)

	// проверяем, что введён один символ
	if len(*delim) != 1 {
		return parametres{}, fmt.Errorf("Delim should be with size 1")
	}
	var filename string
	if len(flag.Args()) == 1 {
		filename = flag.Args()[0]
	}

	params := parametres{
		fields:        columns,
		delim:         *delim,
		withDelimOnly: *isSep,
		filename:      filename,
	}

	return params, nil
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

func doCut(data []string, params parametres, out io.Writer) {
	for i := 0; i < len(data); i++ {
		isDelimExists := strings.Contains(data[i], params.delim)
		if params.withDelimOnly && !isDelimExists {
			continue
		}
		if !params.withDelimOnly && !isDelimExists {
			// fmt.Fprintf(out, "!Conctine %v!\n", isDelimExists)
			fmt.Fprintf(out, "%s\n", data[i])
			continue
		}

		columnStrs := strings.Split(data[i], params.delim)
		for _, colNum := range params.fields {
			if colNum >= len(columnStrs) {
				break
			}
			fmt.Fprintf(out, "%s", columnStrs[colNum])
			if colNum != len(params.fields)-1 {
				fmt.Fprintf(out, "%s", params.delim)
			}
		}
		fmt.Fprint(out, "\n")
	}
}
