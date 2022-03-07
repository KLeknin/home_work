package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			// Place your code here.
		},
		// ...
		// Place your code here.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			_ = tt
		})
	}
}

type tstData struct {
	s, tags string
	res     []ValidationError
	err     error
}

func TestValidateString(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		testData := []tstData{
			{"", "", []ValidationError{}, nil},
			{"", "любые:левые:тэги", []ValidationError{}, nil},
			{"Любая строка", "", []ValidationError{}, nil},
			{"Любая строка", "любые:левые:тэги", []ValidationError{}, nil},
			//len:32 - длина строки должна быть ровно 32 символа;
			{"", "len:0", []ValidationError{}, nil},
			{"1", "len:1", []ValidationError{}, nil},
			{"12345678901234567890123456789012", "len:32", []ValidationError{}, nil},
			//regexp:\\d+ - согласно регулярному выражению строка должна состоять из цифр (\\ - экранирование слэша);
			{"12345678901234567890123456789012", "regexp:\\d+", []ValidationError{}, nil},
			//in:foo,bar - строка должна входить в множество строк {"foo", "bar"}.
			{"a", "in:a", []ValidationError{}, nil},
			{"a", "in:b,a,c", []ValidationError{}, nil},
			{"a", "in:b,c,a", []ValidationError{}, nil},
			{"a", "in:a,b,c", []ValidationError{}, nil},
			{"a", "in:a,a,a", []ValidationError{}, nil},
			//{"", "in:a,a,a", []ValidationError{}, nil},
			//{"", "in:", []ValidationError{}, nil},
			//Допускается комбинация валидаторов по логическому "И" с помощью |, например:
			{"12345678901234567890123456789012", "len:32|regexp:\\d+", []ValidationError{}, nil},
			{"123", "len:3|regexp:\\d+|in:a,123,b", []ValidationError{}, nil},
		}
		for i, testD := range testData {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				//vErr := make(ValidationErrors, 0)
				err, vErr := validateString(testD.s, testD.tags)
				require.Nilf(t, err, "Error: %v", err)
				require.Equal(t, testD.res, vErr, "Results different.\nExpected:%v\nGained:  %v", testD.res, vErr)
				if len(vErr) > 0 {
					println(vErr)
				}
			})
		}
	})
	t.Run("badTag", func(t *testing.T) {
		testData := []tstData{
			{"a", "len:2",
				[]ValidationError{{"a",
					fmt.Errorf("wrong string length, expected:%v", 2),
				}},
				nil},
			{"a", "regexp:\\d+",
				[]ValidationError{{"a",
					fmt.Errorf("string \"%v\" dont match regular tag \"%v\"", "a", "regexp:\\d+"),
				}},
				nil},
			{"a", "in:b,c",
				[]ValidationError{{"a",
					fmt.Errorf("string \"%v\" dont in tag list \"%v\"", "a", "in:b,c"),
				}},
				nil},
		}
		for i, testD := range testData {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				//vErr := make(ValidationErrors, 0)
				err, vErr := validateString(testD.s, testD.tags)
				require.Nilf(t, err, "Error: %v", err)
				for j, _ := range vErr {
					require.Equal(t, testD.res[j].Err.Error(), vErr[j].Err.Error(), "Results different.\nExpected:%v\nGained:  %v", testD.res[j], vErr[j])
				}
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		testData := []tstData{
			{"", "len:", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "len:")},
			{"", "in:", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "in:")},
			{"", "regexp:", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "regexp:")},
			{"", "len", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "len")},
			{"", "in", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "in")},
			{"", "regexp", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "regexp")},
			{"", "regexp:)", []ValidationError{}, fmt.Errorf("tag key dont contain regular expression:\"%v\"", "regexp:)")},
			{"", "len:0|regexp:)", []ValidationError{}, fmt.Errorf("tag key dont contain regular expression:\"%v\"", "regexp:)")},
			{"", "len:стопятьсот", []ValidationError{}, fmt.Errorf("tag key wrong parametr:\"%v\"", "len:стопятьсот")},
		}
		for i, testD := range testData {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				//vErr := make(ValidationErrors, 0)
				err, vErr := validateString(testD.s, testD.tags)
				require.Equal(t, testD.err, err, "No or incorrect error.\nExpected:%v\nGained:  %v", testD.err, err)
				require.Equal(t, testD.res, vErr, "Results different.\nExpected:%v\nGained:  %v", testD.res, vErr)
			})
		}
	})

}
