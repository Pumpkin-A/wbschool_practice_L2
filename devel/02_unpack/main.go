package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const backslashCode = 92

func main() {

}

type codePair struct {
	char   rune
	length int
}

func unpacking(str string) (string, error) {
	var builder strings.Builder

	index := 0
	for index < len(str) {
		pair, nextIndex, err := getCodePair(str, index)
		if err != nil {
			return "", err
		}
		builder.WriteString(strings.Repeat(string(pair.char), pair.length))
		index = nextIndex
	}

	return builder.String(), nil
}

func getCodePair(str string, index int) (codePair, int, error) {
	word, wordSize := utf8.DecodeRuneInString(str[index:])
	if !unicode.IsLetter(word) && int(word) != backslashCode {
		return codePair{}, 0, fmt.Errorf("not expected %d %c", word, word)
	}

	if int(word) == backslashCode {
		slashedWord, slashedWordSize := utf8.DecodeRuneInString(str[index+wordSize:])
		return codePair{char: slashedWord, length: 1}, index + wordSize + slashedWordSize, nil
	}

	indexNumberBegin := index + wordSize
	indexNumberEnd := strings.IndexFunc(str[indexNumberBegin:], func(r rune) bool {
		return !unicode.IsDigit(r)
	})
	if indexNumberEnd == -1 {
		indexNumberEnd = len(str)
	} else {
		indexNumberEnd = indexNumberBegin + indexNumberEnd
	}
	if indexNumberBegin == indexNumberEnd {
		return codePair{char: word, length: 1}, indexNumberEnd, nil
	}

	number, err := strconv.Atoi(str[indexNumberBegin:indexNumberEnd])
	if err != nil {
		return codePair{}, 0, err
	}

	return codePair{char: word, length: number}, indexNumberEnd, nil
}
