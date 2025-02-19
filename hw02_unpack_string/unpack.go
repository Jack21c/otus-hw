package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var res string
	var buffer string
	for _, r := range s {
		// Проверяем, является ли следующий символ цифрой
		if unicode.IsDigit(r) {
			// Числа запрещены
			if buffer == "" {
				return "", ErrInvalidString
			}
			countString := fmt.Sprintf("%c", r)
			count, _ := strconv.Atoi(countString)
			res += strings.Repeat(buffer, count)
			// Очищаем буфер, тк он уже был добавлен к ответу
			buffer = ""
			continue
		}
		// Добавляем значение буфера к ответу
		res += buffer
		// Обновляем буфер текущим символом
		buffer = fmt.Sprintf("%c", r)
	}
	// Добавляем последний символ к ответу
	res += buffer

	return res, nil
}
