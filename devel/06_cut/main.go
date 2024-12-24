package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
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
	params, err := parseArgsIntoParams()
	if err != nil {
		fmt.Println("Invalid arg!", err)
		return
	}

	data, err := readDataFromFile(params.filename)
	if err != nil {
		log.Fatalln(err)
	}

	doCut(data, params, os.Stdout)
}

func parseArgsIntoParams() (parametres, error) {
	fields := flag.String("f", "", "Choose fileds")
	delim := flag.String("d", "", "Choose delim")
	isSep := flag.Bool("s", false, "Separated parametr")

	flag.Parse()

	uniqueColumns := make(map[int]struct{})
	for _, column := range strings.Split(*fields, " ") {
		num, err := strconv.Atoi(column)
		if err != nil {
			return parametres{}, err
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

	if len(*delim) != 1 {
		return parametres{}, fmt.Errorf("Delim should be with size 1")
	}
	if len(flag.Args()) != 1 {
		return parametres{}, fmt.Errorf("Filename should be provided")
	}

	params := parametres{
		fields:        columns,
		delim:         *delim,
		withDelimOnly: *isSep,
		filename:      flag.Args()[0],
	}

	return params, nil
}

func readDataFromFile(name string) ([]string, error) {
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
