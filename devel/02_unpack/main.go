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
	ctxLetter    rune
	ctxLen       int
	ctxBackslash bool
	builder      strings.Builder
}

func newUnpacker() *unpacker {
	return &unpacker{
		builder: strings.Builder{},
	}
}

// naming!!
func (u *unpacker) isCtxLetterEmpty() bool {
	var empty rune
	return u.ctxLetter == empty
}

func (u *unpacker) isCtxBackslash() bool {
	return u.ctxBackslash
}

func (u *unpacker) updateCtxLetter(newChar rune) {
	u.ctxLetter = newChar
}

func (u *unpacker) updateCtxLen(newLen int) {
	u.ctxLen = u.ctxLen*10 + newLen
}

func (u *unpacker) updateCtxBackslash() {
	u.ctxBackslash = true
}

func (u *unpacker) resetCtx() {
	var empty rune
	u.ctxLetter = empty
	u.ctxLen = 0
	u.ctxBackslash = false
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
	} else if u.isCtxBackslash() {
		return fmt.Errorf("incorrect string: unused backslash")
	}
	u.updateCtxLetter(r)
	return nil
}

func (u *unpacker) digitsHandler(len rune) error {
	// если контекст пуст и не включен контекст слэша - ошибка (нет символов для повтора)
	// если контекст слэша включен, значит цифра - символ для обработки
	if u.isCtxLetterEmpty() {
		if u.isCtxBackslash() {
			u.updateCtxLetter(len)
			return nil
		}
		return fmt.Errorf("incorrect string")
	}
	// встретили очередную цифру в строке - обновляем длину (это длина символа, который сейчас в контексте)
	u.updateCtxLen(int(len - '0'))
	return nil
}

func (u *unpacker) backslashHandler() error {
	// если контекст слэша включен и символа для обработки нет - слэш и есть символ для обработки
	// если контекст слэша включен и контекст символа заполнен - значит записываем накопленный контекст в итоговую строку
	if u.isCtxBackslash() {
		if u.isCtxLetterEmpty() {
			u.updateCtxLetter(rune(backslashCode))
			return nil
		}
		u.writeCtxToString()
		u.updateCtxBackslash()
		return nil
	}

	if !u.isCtxLetterEmpty() {
		u.writeCtxToString()
	}
	u.updateCtxBackslash()
	return nil
}

func unpacking(str string) (string, error) {
	fmt.Println(str)
	unpacker := newUnpacker()
	runes := []rune(str)

	for _, v := range runes {
		if unicode.IsDigit(v) {
			// обработчик, когда встречаем число
			err := unpacker.digitsHandler(v)
			if err != nil {
				return "", err
			}
			continue
		} else if int(v) == backslashCode {
			// обработчик после слэша
			err := unpacker.backslashHandler()
			if err != nil {
				return "", err
			}
			continue
		}

		err := unpacker.lettersHandler(v)
		if err != nil {
			return "", err
		}

	}

	//обработка накопленного контекста после завершения строки
	if !unpacker.isCtxLetterEmpty() {
		unpacker.writeCtxToString()
	} else if unpacker.ctxBackslash {
		return "", fmt.Errorf("incorrect string: unused backslash")
	}

	fmt.Println(unpacker.builder.String())
	return unpacker.builder.String(), nil
}
