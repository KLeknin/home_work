package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inStr string) (string, error) {
	var retStr strings.Builder
	var retErr error
	dblSlash := false

	if inStr != "" {
		var lastRune rune
		for _, aRune := range inStr + "e" { //добавляем один "лишний" символ для удобства, не надо будет обрабатывать "хвостик" за пределами цикла
			aRune := aRune
			isDight := unicode.IsDigit(aRune)
			switch {
			case isDight:
				{
					switch {
					case unicode.IsDigit(lastRune):
						{ //две цифры подряд
							retErr = ErrInvalidString
						}
					case (lastRune == 92):
						{ // `qwe\4\5`
							lastRune = aRune
							break
						}
					case lastRune != 0:
						{ //"a2b3"

							/* по "подсказке" надо так:
							runeCount, err := strconv.Atoi(string(aRune))
							if err != nil {
								retErr = ErrInvalidString
							}
								но тут проще без Atoi сделать: */
							runeCount := int(aRune - 48)
							for i := 0; i < runeCount; i++ {
								retStr.WriteRune(lastRune)
							}
							lastRune = 0

						}
					case dblSlash:
						{ // "\\3\\5"
							runeCount := int(aRune-48) - 1 //первый слеш уже добавлен
							if runeCount >= 0 {
								for i := 0; i < runeCount; i++ {
									retStr.WriteRune(92)
								}
							} else {

								//Есть еще такое редкое извращение "\\0"

								s := retStr.String()
								s = s[:len(s)-1] //тут можно, т.к. руна "\" занимает 1 байт
								retStr.Reset()
								retStr.WriteString(s)
							}

							dblSlash = false
						}

					default:
						{
							retErr = ErrInvalidString
						}
					}
				}
			case (lastRune == 92) && (aRune == 92): // `\\`
				{
					retStr.WriteRune(lastRune)
					lastRune = 0
					dblSlash = true
				}
			case (lastRune == 92):
				{
					retErr = ErrInvalidString
				}
			case lastRune != 0: //просто буква
				{
					retStr.WriteRune(lastRune)
					lastRune = aRune
				}
			default:
				{
					lastRune = aRune
					dblSlash = false
				}
			}
			if retErr != nil {
				break
			}
		}
	}
	if retErr != nil {
		retStr.Reset()
	}
	return retStr.String(), retErr
}
