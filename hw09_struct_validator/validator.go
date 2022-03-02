package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func validateString(s, tag string) error {
	fmt.Printf("%v, %T\n\n", s, s)
	return nil
}
func validateInt(i int64, tag string) error {
	fmt.Printf("%v, %T\n\n", i, i)
	return nil
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
						err = validateString(fv.String(), tagStr)
					case reflect.Int:
						err = validateInt(fv.Int(), tagStr)
					}
					//todo обработать err
				}
			}
		case reflect.String:
			err = validateString(fieldValue.String(), tagStr)
			//todo обработать err
		case reflect.Int:
			err = validateInt(fieldValue.Int(), tagStr)
			//todo обработать err
		}
	}
	_ = err
	return nil
}
