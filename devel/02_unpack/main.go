package main

import (
	"fmt"
	"strings"
	"unicode"
)

const backslashCode = 92

type Unpacker struct {
	//контекст
	Symbol    rune
	Len       int
	Backslash bool

	builder strings.Builder
}

func NewUnpacker() *Unpacker {
	return &Unpacker{
		builder: strings.Builder{},
	}
}

func (u *Unpacker) IsSymbolEmpty() bool {
	return u.Symbol == rune(0)
}

func (u *Unpacker) IsBackslash() bool {
	return u.Backslash
}

func (u *Unpacker) UpdateSymbol(newChar rune) {
	u.Symbol = newChar
}

func (u *Unpacker) UpdateLen(newLen int) {
	u.Len = u.Len*10 + newLen
}

func (u *Unpacker) UpdateBackslash() {
	u.Backslash = true
}

func (u *Unpacker) ResetCtx() {
	u.Symbol = rune(0)
	u.Len = 0
	u.Backslash = false
}

// записывает контекст в итоговую строку указанное количество раз
func (u *Unpacker) WriteCtxToString() {
	var len int
	if u.Len == 0 {
		len = 1
	} else {
		len = u.Len
	}
	u.builder.WriteString(strings.Repeat(string(u.Symbol), len))
	// очистка контекста после записи
	u.ResetCtx()
}

func (u *Unpacker) SymbolsHandler(r rune) error {
	// встретили новый символ буквы. Если контекст не пуст - надо записать данные в итоговую строку
	// После обработки записываем в контекст новый символ
	if !u.IsSymbolEmpty() {
		u.WriteCtxToString()
	} else if u.IsBackslash() {
		return fmt.Errorf("incorrect string: unused backslash")
	}
	u.UpdateSymbol(r)
	return nil
}

func (u *Unpacker) DigitsHandler(len rune) error {
	// если контекст пуст и не включен контекст слэша - ошибка (нет символов для повтора)
	// если контекст слэша включен, значит цифра - символ для обработки
	if u.IsSymbolEmpty() {
		if u.IsBackslash() {
			u.UpdateSymbol(len)
			return nil
		}
		return fmt.Errorf("incorrect string")
	}
	// встретили очередную цифру в строке - обновляем длину (это длина символа, который сейчас в контексте)
	u.UpdateLen(int(len - '0'))
	return nil
}

func (u *Unpacker) BackslashHandler() error {
	// если контекст слэша включен и символа для обработки нет - слэш и есть символ для обработки
	// если контекст слэша включен и контекст символа заполнен - значит записываем накопленный контекст в итоговую строку
	if u.IsBackslash() {
		if u.IsSymbolEmpty() {
			u.UpdateSymbol(rune(backslashCode))
			return nil
		}
		u.WriteCtxToString()
		u.UpdateBackslash()
		return nil
	}

	if !u.IsSymbolEmpty() {
		u.WriteCtxToString()
	}
	u.UpdateBackslash()
	return nil
}

func (u *Unpacker) Do(str string) (string, error) {
	for _, v := range str {
		if unicode.IsDigit(v) {
			// обработчик, когда встречаем число
			err := u.DigitsHandler(v)
			if err != nil {
				return "", err
			}
			continue
		} else if int(v) == backslashCode {
			// обработчик после слэша
			err := u.BackslashHandler()
			if err != nil {
				return "", err
			}
			continue
		}

		err := u.SymbolsHandler(v)
		if err != nil {
			return "", err
		}
	}

	//обработка накопленного контекста после завершения строки
	if !u.IsSymbolEmpty() {
		u.WriteCtxToString()
	} else if u.Backslash {
		return "", fmt.Errorf("incorrect string: unused backslash")
	}

	fmt.Println(u.builder.String())
	return u.builder.String(), nil
}

func main() {

}
