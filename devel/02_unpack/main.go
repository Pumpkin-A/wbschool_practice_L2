package main

import (
	"fmt"
	"strings"
	"unicode"
)

const backslashCode = 92

func main() {

}

type unpacker struct {
	ctxLetter rune
	ctxLen    int
	slash     bool
	builder   strings.Builder
}

func newUnpacker() *unpacker {
	return &unpacker{
		builder: strings.Builder{},
	}
}

func (u *unpacker) isCtxLetterEmpty() bool {
	var empty rune
	return u.ctxLetter == empty
}

func (u *unpacker) updateCtxLetter(newChar rune) {
	u.ctxLetter = newChar
}

func (u *unpacker) updateCtxLen(newLen int) {
	u.ctxLen = u.ctxLen*10 + newLen
}

func (u *unpacker) resetCtx() {
	var empty rune
	u.ctxLetter = empty
	u.ctxLen = 0
	u.slash = false
}

// записывает контекст в итоговую строку указанное количество раз
func (u *unpacker) writeCtxToString() {
	var len int
	if u.ctxLen == 0 {
		len = 1
	} else {
		len = u.ctxLen
	}
	u.builder.WriteString(strings.Repeat(string(u.ctxLetter), len))
	// очистка контекста после записи
	u.resetCtx()
}

func (u *unpacker) lettersHandler(r rune) error {
	// встретили новый символ буквы. Если контекст не пуст - надо записать данные в итоговую строку
	// После обработки записываем в контекст новый символ
	if !u.isCtxLetterEmpty() {
		u.writeCtxToString()
	}
	u.updateCtxLetter(r)
	return nil
}

func (u *unpacker) digitsHandler(len int) error {
	// если контекст пуст - ошибка (строка начинается с цифры - некорректная строка)
	if u.isCtxLetterEmpty() {
		return fmt.Errorf("incorrect string")
	}
	// встретили очередную цифру в строке - обновляем длину (это длина символа, который сейчас в контексте)
	u.updateCtxLen(len)
	return nil
}

func unpacking(str string) (string, error) {
	fmt.Println(str)
	unpacker := newUnpacker()
	runes := []rune(str)

	for _, v := range runes {
		if unicode.IsDigit(v) {
			// обработчик, когда встречаем число
			err := unpacker.digitsHandler(int(v - '0'))
			if err != nil {
				return "", err
			}
			continue
		}
		// else if int(v) == backslashCode {
		// 	// обработчик после слэша
		// 	continue
		// }

		err := unpacker.lettersHandler(v)
		if err != nil {
			return "", err
		}
		// обработка для букв и иных символов
	}

	if !unpacker.isCtxLetterEmpty() {
		unpacker.writeCtxToString()
	}
	fmt.Println(unpacker.builder.String())
	return unpacker.builder.String(), nil
}

// func unpacking(str string) (string, error) {
// 	var builder strings.Builder

// 	index := 0
// 	runes := []rune(str)
// 	for index < len(runes) {
// 		pair, nextIndex, err := getCodePair(runes, index)
// 		if err != nil {
// 			return "", err
// 		}
// 		builder.WriteString(strings.Repeat(string(pair.char), pair.length))
// 		fmt.Println(builder.String())
// 		index = nextIndex
// 	}

// 	return builder.String(), nil
// }

// getPair {
// word, wordSize := utf8.DecodeRuneInString(str[index:])
// if !unicode.IsLetter(word) && int(word) != backslashCode {
// 	return codePair{}, 0, fmt.Errorf("not expected %d %c", word, word)
// }

// if int(word) == backslashCode {
// 	slashedWord, slashedWordSize := utf8.DecodeRuneInString(str[index+wordSize:])
// 	return codePair{char: slashedWord, length: 1}, index + wordSize + slashedWordSize, nil
// }

// indexNumberBegin := index + wordSize
// indexNumberEnd := strings.IndexFunc(str[indexNumberBegin:], func(r rune) bool {
// 	return !unicode.IsDigit(r)
// })
// if indexNumberEnd == -1 {
// 	indexNumberEnd = len(str)
// } else {
// 	indexNumberEnd = indexNumberBegin + indexNumberEnd
// }
// if indexNumberBegin == indexNumberEnd {
// 	return codePair{char: word, length: 1}, indexNumberEnd, nil
// }

// number, err := strconv.Atoi(str[indexNumberBegin:indexNumberEnd])
// if err != nil {
// 	return codePair{}, 0, err
// }

// return codePair{char: word, length: number}, indexNumberEnd, nil
// }
