package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v *ValidationError) Error() string { return v.Field + ": " + v.Err.Error() }

type ValidationErrors []ValidationError

func (v *ValidationErrors) Error() string {
	s := ""
	for _, e := range *v {
		s += e.Error() + "\n"
	}
	return s
}

func validateString(s, tags string) (vErr []ValidationError, err error) { //nolint: gocognit
	vErr = make([]ValidationError, 0)

	for _, oneTag := range strings.Split(tags, "|") {
		if oneTag == "" {
			continue
		}
		tagKeyValue := strings.SplitN(oneTag, ":", 2)
		if len(tagKeyValue) < 2 || tagKeyValue[1] == "" {
			return vErr, fmt.Errorf("tag key without value:\"%v\"", oneTag)
		}
		switch tagKeyValue[0] {
		case "len": // len:32 - длина строки должна быть ровно 32 символа;
			{
				lenExpected, err := strconv.Atoi(strings.TrimSpace(tagKeyValue[1]))
				if err != nil {
					return vErr, fmt.Errorf("tag key wrong parametr:\"%v\"", oneTag)
				}
				if lenExpected != len(s) {
					vErr = append(vErr, ValidationError{s, fmt.Errorf("wrong string length, expected:%v", lenExpected)})
				}
			}

		case "regexp": // regexp:\\d+ - согласно регулярному выражению строка должна состоять из цифр
			{
				rExp, err := regexp.Compile(tagKeyValue[1])
				if err != nil {
					return vErr, fmt.Errorf("tag key dont contain regular expression:\"%v\"", oneTag)
				}
				if !rExp.MatchString(s) {
					vErr = append(vErr, ValidationError{s, fmt.Errorf("string \"%v\" dont match regular tag \"%v\"", s, oneTag)})
				}
			}

		case "in": //in:foo,bar - строка должна входить в множество строк {"foo", "bar"}.
			{
				found := false
				for _, str := range strings.Split(tagKeyValue[1], ",") {
					if str == s {
						found = true
						break
					}
				}
				if !found {
					vErr = append(vErr, ValidationError{s, fmt.Errorf("string \"%v\" dont in tag list \"%v\"", s, oneTag)})
				}
			}
		}
	}
	return vErr, nil
}

func validateInt(i int64, tags string) (vErr []ValidationError, err error) { //nolint: gocognit
	vErr = make([]ValidationError, 0)

	for _, oneTag := range strings.Split(tags, "|") {
		if oneTag == "" {
			continue
		}
		tagKeyValue := strings.SplitN(oneTag, ":", 2)
		if len(tagKeyValue) < 2 || tagKeyValue[1] == "" {
			return vErr, fmt.Errorf("tag key without value:\"%v\"", oneTag)
		}
		switch tagKeyValue[0] {
		case "min": // min:10 - число не может быть меньше 10;
			{
				min, err := strconv.Atoi(strings.TrimSpace(tagKeyValue[1]))
				if err != nil {
					return vErr, fmt.Errorf("tag key wrong parametr:\"%v\"", oneTag)
				}
				if i < int64(min) {
					vErr = append(vErr, ValidationError{
						strconv.Itoa(int(i)),
						fmt.Errorf("value too small: %v < %v, %v", i, min, oneTag),
					})
				}
			}

		case "max": // max:20 - число не может быть больше 20;
			{
				max, err := strconv.Atoi(strings.TrimSpace(tagKeyValue[1]))
				if err != nil {
					return vErr, fmt.Errorf("tag key wrong parametr:\"%v\"", oneTag)
				}
				if i > int64(max) {
					vErr = append(vErr, ValidationError{
						strconv.Itoa(int(i)),
						fmt.Errorf("value too big: %v > %v, %v", i, max, oneTag),
					})
				}
			}

		case "in": // in:256,1024 - число должно входить в множество чисел {256, 1024};
			{
				found := false
				for _, str := range strings.Split(tagKeyValue[1], ",") {
					j, err := strconv.Atoi(strings.TrimSpace(str))
					if err != nil {
						return vErr,
							fmt.Errorf("not integer value in parametr:\"%v\", error:%v", oneTag, err.Error()) //nolint: errorlint
					}
					if int64(j) == i {
						found = true
						break
					}
				}
				if !found {
					vErr = append(vErr, ValidationError{
						strconv.Itoa(int(i)),
						fmt.Errorf("integer \"%v\" dont in tag list \"%v\"", i, oneTag),
					})
				}
			}
		}
	}

	return vErr, nil
}

func Validate(v interface{}) error {
	if v == nil {
		return fmt.Errorf("data is nil: %v", v)
	}
	vr := reflect.ValueOf(v)
	if vr.Kind() != reflect.Struct {
		return fmt.Errorf("data is not structure: %v", v)
	}
	var validationErrors ValidationErrors
	var vErr []ValidationError
	var err error
	st := reflect.ValueOf(v)
	for i := 0; i < st.NumField(); i++ {
		fieldValue := st.Field(i)
		fieldType := st.Type().Field(i)
		tagStr, ok := fieldType.Tag.Lookup("validate")
		if !ok {
			continue
		}

		switch fieldValue.Kind() { //nolint: exhaustive
		case reflect.Slice:
			{
				for j := 0; j < fieldValue.Len(); j++ {
					fv := fieldValue.Index(j)
					switch fv.Kind() { //nolint: exhaustive
					case reflect.String:
						vErr, err = validateString(fv.String(), tagStr)
					case reflect.Int:
						vErr, err = validateInt(fv.Int(), tagStr)
					}
					if err != nil {
						return fmt.Errorf("error in field %v, element %v, tag \"%v\": %v",
							fieldValue.String(), j, tagStr, err.Error()) //nolint: errorlint
					}
					validationErrors = append(validationErrors, vErr...)
					vErr = []ValidationError{}
				}
			}
		case reflect.String:
			vErr, err = validateString(fieldValue.String(), tagStr)
		case reflect.Int:
			vErr, err = validateInt(fieldValue.Int(), tagStr)
		}
		if err != nil {
			return fmt.Errorf("error in field %v, tag \"%v\": %v", fieldValue.String(), tagStr, err.Error()) //nolint: errorlint
		}
		validationErrors = append(validationErrors, vErr...)
		vErr = []ValidationError{}
	}
	if len(validationErrors) > 0 {
		return &validationErrors
	}
	return nil
}
