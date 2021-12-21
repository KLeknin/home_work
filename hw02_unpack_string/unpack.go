package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func unpackRuneList(inRunes []rune) (string, error) {
	var retStr strings.Builder
	var retErr error
	var inLen, pos int
	var s string
	var n int
	inLen = len(inRunes)
	retStr.Reset()
	for pos < inLen {
		switch {
		case inRunes[pos] == 92: // слеш
			{
				s, n, retErr = unSlash(inRunes[pos:])
				if retErr != nil {
					return "", ErrInvalidString
				}
				pos += n
				retStr.WriteString(s)
			}
		case unicode.IsDigit(inRunes[pos]): // цифра
			{
				return "", ErrInvalidString
			}
		default:
			{
				if len(inRunes) > pos+1 {
					s, n = needDight(inRunes[pos+1:], inRunes[pos]) // ищем множитель
					pos += n
					retStr.WriteString(s)
				} else {
					retStr.WriteRune(inRunes[pos])
				}
			}
		}
		pos++
	}
	return retStr.String(), nil
}

func unSlash(sRunes []rune) (string, int, error) {
	if (len(sRunes) > 1) && ((sRunes[1] == 92) || (unicode.IsDigit(sRunes[1]))) { // двойной слеш или цифра после слеша
		if len(sRunes) >= 2 {
			s, n := needDight(sRunes[2:], sRunes[1]) // ищем множитель
			return s, 1 + n, nil
		}
		return string(sRunes[1]), 1, nil
	}
	return "", 0, ErrInvalidString
}

func needDight(sRunes []rune, aRune rune) (string, int) {
	var s strings.Builder
	s.Reset()
	n := 1
	k := 0
	if (len(sRunes) >= 1) && unicode.IsDigit(sRunes[0]) {
		n = int(sRunes[0] - 48) // Для одной цифры проще так, чем через strconv.Atoi
		k = 1
	}
	for i := 0; i < n; i++ {
		s.WriteRune(aRune)
	}
	return s.String(), k
}

func Unpack(inStr string) (string, error) {
	retStr, retErr := unpackRuneList([]rune(inStr))
	return retStr, retErr
}
