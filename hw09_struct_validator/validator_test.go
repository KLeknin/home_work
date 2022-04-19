package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
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
			App{""},
			&ValidationErrors{{"", fmt.Errorf("wrong string length, expected:%v", 5)}},
		},
		{App{"12345"}, nil},
		{
			App{"12345678"},
			&ValidationErrors{{"12345678", fmt.Errorf("wrong string length, expected:%v", 5)}},
		},
		{Token{[]byte{1, 2, 3, 4}, []byte{1, 2, 3, 4}, []byte{1, 2, 3, 4}}, nil},
		{Response{200, ""}, nil},
		{
			Response{202, ""},
			&ValidationErrors{{"202", fmt.Errorf("integer \"%v\" dont in tag list \"%v\"", "202", "in:200,404,500")}},
		},
		{User{
			"123456789012345678901234567890123456", "Peter", 18, "No@mail.me", "stuff",
			[]string{"12345678901", "11234567890"},
			[]byte{0},
		}, nil},
		{
			User{
				"123", "Peter", 16, "Not e-mail me", "friend",
				[]string{"12345678901", "12345678"},
				[]byte{0},
			},
			&ValidationErrors{
				{
					"123",
					fmt.Errorf("wrong string length, expected:%v", 36),
				},
				{
					"16",
					fmt.Errorf("value too small: %v < %v, %v", 16, 18, "min:18"),
				},
				{
					"Not e-mail me",
					fmt.Errorf("string \"%v\" dont match regular tag \"%v\"", "Not e-mail me", "regexp:^\\w+@\\w+\\.\\w+$"),
				},
				{
					"friend",
					fmt.Errorf("string \"%v\" dont in tag list \"%v\"", "friend", "in:admin,stuff"),
				},
				{
					"12345678",
					errors.New("wrong string length, expected:11"),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			// Place your code here.
			gainedErr := Validate(tt.in)
			require.Equal(t, tt.expectedErr, gainedErr,
				"Results different.\nExpected:%v\nGained:  %v", tt.expectedErr, gainedErr)
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
			// len:32 - длина строки должна быть ровно 32 символа;
			{"", "len:0", []ValidationError{}, nil},
			{"1", "len:1", []ValidationError{}, nil},
			{"12345678901234567890123456789012", "len:32", []ValidationError{}, nil},
			// regexp:\\d+ - согласно регулярному выражению строка должна состоять из цифр (\\ - экранирование слэша);
			{"12345678901234567890123456789012", "regexp:\\d+", []ValidationError{}, nil},
			//in:foo,bar - строка должна входить в множество строк {"foo", "bar"}.
			{"a", "in:a", []ValidationError{}, nil},
			{"a", "in:b,a,c", []ValidationError{}, nil},
			{"a", "in:b,c,a", []ValidationError{}, nil},
			{"a", "in:a,b,c", []ValidationError{}, nil},
			{"a", "in:a,a,a", []ValidationError{}, nil},
			// Допускается комбинация валидаторов по логическому "И" с помощью |, например:
			{"12345678901234567890123456789012", "len:32|regexp:\\d+", []ValidationError{}, nil},
			{"123", "len:3|regexp:\\d+|in:a,123,b", []ValidationError{}, nil},
		}
		for i, testD := range testData {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				vErr, err := validateString(testD.s, testD.tags)
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
			{
				"a", "len:2",
				[]ValidationError{{"a", fmt.Errorf("wrong string length, expected:%v", 2)}},
				nil,
			},
			{
				"a", "regexp:\\d+",
				[]ValidationError{{"a", fmt.Errorf("string \"%v\" dont match regular tag \"%v\"", "a", "regexp:\\d+")}},
				nil,
			},
			{
				"a", "in:b,c",
				[]ValidationError{{"a", fmt.Errorf("string \"%v\" dont in tag list \"%v\"", "a", "in:b,c")}},
				nil,
			},
		}
		for i, testD := range testData {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				// vErr := make(ValidationErrors, 0)
				vErr, err := validateString(testD.s, testD.tags)
				require.Nilf(t, err, "Error: %v", err)
				for j := range vErr {
					require.Equal(t, testD.res[j].Err.Error(), vErr[j].Err.Error(),
						"Results different.\nExpected:%v\nGained:  %v", testD.res[j], vErr[j])
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
			{
				"", "regexp:)",
				[]ValidationError{},
				fmt.Errorf("tag key dont contain regular expression:\"%v\"", "regexp:)"),
			},
			{
				"", "len:0|regexp:)",
				[]ValidationError{},
				fmt.Errorf("tag key dont contain regular expression:\"%v\"", "regexp:)"),
			},
			{"", "len:стопятьсот", []ValidationError{}, fmt.Errorf("tag key wrong parametr:\"%v\"", "len:стопятьсот")},
		}
		for i, testD := range testData {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				// vErr := make(ValidationErrors, 0)
				vErr, err := validateString(testD.s, testD.tags)
				require.Equal(t, testD.err, err, "No or incorrect error.\nExpected:%v\nGained:  %v", testD.err, err)
				require.Equal(t, testD.res, vErr, "Results different.\nExpected:%v\nGained:  %v", testD.res, vErr)
			})
		}
	})
}

type tstDataInt struct {
	i    int64
	tags string
	res  []ValidationError
	err  error
}

func TestValidateInt(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		testData := []tstDataInt{
			{0, "", []ValidationError{}, nil},
			// min:10 - число не может быть меньше 10;
			{10, "min:10", []ValidationError{}, nil},
			{11, "min:10", []ValidationError{}, nil},
			// max:20 - число не может быть больше 20;
			{20, "max:20", []ValidationError{}, nil},
			{19, "max:20", []ValidationError{}, nil},
			// in:256,1024 - число должно входить в множество чисел {256, 1024};
			{1, "in:1,2,3", []ValidationError{}, nil},
			{2, "in:1,2,3", []ValidationError{}, nil},
			{3, "in: 1 , 2 , 3 ", []ValidationError{}, nil},
		}
		for j, testD := range testData {
			t.Run(fmt.Sprint(j), func(t *testing.T) {
				vErr, err := validateInt(testD.i, testD.tags)
				require.Nilf(t, err, "Error: %v", err)
				require.Equal(t, testD.res, vErr, "Results different.\nExpected:%v\nGained:  %v", testD.res, vErr)
				if len(vErr) > 0 {
					println(vErr)
				}
			})
		}
	})
	t.Run("badTag", func(t *testing.T) {
		testData := []tstDataInt{
			{
				1, "in:2,3",
				[]ValidationError{{"1", fmt.Errorf("integer \"%v\" dont in tag list \"%v\"", 1, "in:2,3")}},
				nil,
			},
			{
				9, "min:10",
				[]ValidationError{{strconv.Itoa(9), fmt.Errorf("value too small: %v < %v, %v", 9, 10, "min:10")}},
				nil,
			},
			{
				21, "max:20",
				[]ValidationError{{strconv.Itoa(21), fmt.Errorf("value too big: %v > %v, %v", 21, 20, "max:20")}},
				nil,
			},
		}
		for j, testD := range testData {
			t.Run(fmt.Sprint(j), func(t *testing.T) {
				vErr, err := validateInt(testD.i, testD.tags)
				require.Nilf(t, err, "Error: %v", err)
				for j := range vErr {
					require.Equal(t, testD.res[j].Err.Error(), vErr[j].Err.Error(),
						"Results different.\nExpected:%v\nGained:  %v", testD.res[j], vErr[j])
				}
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		_, err := strconv.Atoi(")")
		testData := []tstDataInt{
			{0, "min:", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "min:")},
			{0, "max:", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "max:")},
			{0, "in:", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "in:")},
			{0, "min", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "min")},
			{0, "max", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "max")},
			{0, "in", []ValidationError{}, fmt.Errorf("tag key without value:\"%v\"", "in")},
			{0, "min:стопятьсот", []ValidationError{}, fmt.Errorf("tag key wrong parametr:\"%v\"", "min:стопятьсот")},
			{0, "max:стопятьсот", []ValidationError{}, fmt.Errorf("tag key wrong parametr:\"%v\"", "max:стопятьсот")},
			{0, "in:),0", []ValidationError{}, fmt.Errorf("not integer value in parametr:\"%v\", error:%v", "in:),0",
				err)}, //nolint: errorlint
		}
		for i, testD := range testData {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				vErr, err := validateInt(testD.i, testD.tags)
				require.Equal(t, testD.err, err, "No or incorrect error.\nExpected:%v\nGained:  %v", testD.err, err)
				require.Equal(t, testD.res, vErr, "Results different.\nExpected:%v\nGained:  %v", testD.res, vErr)
			})
		}
	})
}
