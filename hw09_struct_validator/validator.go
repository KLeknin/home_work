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

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func validateString(s, tags string) (err error, vErr []ValidationError) {
	vErr = make([]ValidationError, 0, 0)
	tagKeyValue := make([]string, 0, 2)
	for _, oneTag := range strings.Split(tags, "|") {
		if oneTag == "" {
			continue
		}
		tagKeyValue = strings.SplitN(oneTag, ":", 2)
		if len(tagKeyValue) < 2 || tagKeyValue[1] == "" {
			return fmt.Errorf("tag key without value:\"%v\"", oneTag), vErr
		}
		switch tagKeyValue[0] {
		case "len": //len:32 - длина строки должна быть ровно 32 символа;
			{
				lenExpected, err := strconv.Atoi(tagKeyValue[1])
				if err != nil {
					return fmt.Errorf("tag key wrong parametr:\"%v\"", oneTag), vErr
				}
				if lenExpected != len(s) {
					vErr = append(vErr, ValidationError{s, fmt.Errorf("wrong string length, expected:%v", lenExpected)})
				}
			}
		case "regexp": //regexp:\\d+ - согласно регулярному выражению строка должна состоять из цифр (\\ - экранирование слэша);
			{
				rExp, err := regexp.Compile(tagKeyValue[1])
				if err != nil {
					return fmt.Errorf("tag key dont contain regular expression:\"%v\"", oneTag), vErr
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
	return nil, vErr
}
func validateInt(i int64, tag string) (err error, vErr ValidationErrors) {
	vErr = make([]ValidationError, 0, 0)
	//tagKeyValue := make([]string, 0, 2)
	return nil, vErr
}

func Validate(v interface{}) error {
	// Place your code here.
	if v == nil { //todo test
		return fmt.Errorf("data is nil: %v", v)
	}
	vr := reflect.ValueOf(v)
	if vr.Kind() != reflect.Struct { //todo test
		return fmt.Errorf("data is not structure: %v", v)
	}

	validationErrors := make([]ValidationError, 0)
	vErr := make([]ValidationError, 0)
	var err error
	st := reflect.ValueOf(v)
	for i := 0; i < st.NumField(); i++ {
		fieldValue := st.Field(i)
		fieldType := st.Type().Field(i)
		tagStr, ok := fieldType.Tag.Lookup("validate")
		if !ok {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Slice:
			{
				for i := 0; i < fieldValue.Len(); i++ {
					fv := fieldValue.Index(i)
					switch fv.Kind() {
					case reflect.String:
						err, vErr = validateString(fv.String(), tagStr)
					case reflect.Int:
						err, vErr = validateInt(fv.Int(), tagStr)
					}
					//todo обработать err
					validationErrors = append(validationErrors, vErr...)
				}
			}
		case reflect.String:
			err, vErr = validateString(fieldValue.String(), tagStr)
			//todo обработать err
		case reflect.Int:
			err, vErr = validateInt(fieldValue.Int(), tagStr)
			//todo обработать err
		}
	}
	_ = err
	return nil
}
