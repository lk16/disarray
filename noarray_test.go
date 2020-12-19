package noarray

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type foo struct {
	Zero string  `json:"offset_zero"`
	One  float64 `json:"offset_one"`
	Two  string  `json:"offset_two"`
}

func TestNoArraySimple(t *testing.T) {
	type testCase struct {
		name           string
		input          []byte
		expectedOutput foo
		expectedErr    error
	}

	testCases := []testCase{
		{
			name:           "baseCase",
			input:          []byte(`["foo",123.4,"bar"]`),
			expectedOutput: foo{Zero: "foo", One: 123.4, Two: "bar"},
			expectedErr:    nil,
		},
		{
			name:           "arrayTooShort",
			input:          []byte(`["foo",123.4]`),
			expectedOutput: foo{Zero: "foo", One: 123.4, Two: ""},
			expectedErr:    nil,
		},
		{
			name:           "wrongFieldType",
			input:          []byte(`["foo","baz","bar"]`),
			expectedOutput: foo{Zero: "foo", Two: "bar"},
			expectedErr:    errors.New("json: cannot unmarshal string into Go struct field foo.offset_one of type float64"),
		},
		{
			name:           "brokenJSON",
			input:          []byte(`["`),
			expectedOutput: foo{},
			expectedErr:    errors.New("unexpected end of JSON input"),
		},
		{
			name:           "notASlice",
			input:          []byte(`{}`),
			expectedOutput: foo{},
			expectedErr:    errors.New("json: cannot unmarshal object into Go value of type []interface {}"),
		},
		{
			name:           "extraSliceValues",
			input:          []byte(`["foo",123.4,"bar","somethingextra",null,123,true]`),
			expectedOutput: foo{Zero: "foo", One: 123.4, Two: "bar"},
			expectedErr:    nil,
		},
		{
			name:           "arrayTooLong",
			input:          []byte(`["foo",123.4,"bar",1,1,1,1,1,1,1,1,1,1]`),
			expectedOutput: foo{},
			expectedErr:    ErrArrayTooLong},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var output foo
			err := UnmarshalAsObject(testCase.input, &output)

			if testCase.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, testCase.expectedErr.Error(), err.Error())
			}

			assert.Equal(t, testCase.expectedOutput, output)
		})
	}
}

type baz struct {
	Key string `json:"key"`
}

type bar struct {
	Zero  string         `json:"offset_zero"`
	One   []baz          `json:"offset_one"`
	Two   string         `json:"offset_two"`
	Three map[string]baz `json:"offset_three"`
	Four  float64        `json:"offset_four"`
}

func TestNoArrayAdvanced(t *testing.T) {
	type testCase struct {
		name           string
		input          []byte
		expectedOutput bar
		expectedErr    error
	}

	testCases := []testCase{
		{
			name:  "baseCase",
			input: []byte(`["something",[{"key": "value"}, {"key": "othervalue"}], "two", {"asdf": {"key": "somevalue"}}]`),
			expectedOutput: bar{
				Zero: "something",
				One: []baz{
					{Key: "value"},
					{Key: "othervalue"},
				},
				Two: "two",
				Three: map[string]baz{
					"asdf": {Key: "somevalue"},
				},
			},
			expectedErr: nil,
		},
		{
			name:  "baseCase",
			input: []byte(`["",[], "", {}]`),
			expectedOutput: bar{
				One:   []baz{},
				Three: map[string]baz{},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var output bar
			err := UnmarshalAsObject(testCase.input, &output)

			if testCase.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, testCase.expectedErr.Error(), err.Error())
			}

			assert.Equal(t, testCase.expectedOutput, output)
		})
	}
}