package disarray

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type foo struct {
	Zero string  `json:"0"`
	One  float64 `json:"1"`
	Two  string  `json:"2"`
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
			expectedOutput: foo{Zero: "foo"},
			expectedErr:    errors.New("json: cannot unmarshal string into Go value of type float64"),
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
			expectedErr:    errors.New("json: cannot unmarshal object into Go value of type []json.RawMessage"),
		},
		{
			name:           "extraSliceValues",
			input:          []byte(`["foo",123.4,"bar","somethingextra",null,123,true]`),
			expectedOutput: foo{Zero: "foo", One: 123.4, Two: "bar"},
			expectedErr:    nil,
		},
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
	Zero  string         `json:"0"`
	One   []baz          `json:"1"`
	Two   string         `json:"2"`
	Three map[string]baz `json:"3"`
	Four  float64        `json:"4"`
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
			name:  "advanced",
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
			name:  "empty",
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

func TestBrokenTag(t *testing.T) {

	type brokenTag struct {
		Field string `json:"noNumberHere"`
	}

	bytes := []byte(`["value"]`)

	var broken brokenTag
	err := UnmarshalAsObject(bytes, &broken)

	assert.Equal(t, err.Error(), `strconv.Atoi: parsing "noNumberHere": invalid syntax`)
}

func TestOtherTag(t *testing.T) {

	type otherTag struct {
		Field string `xml:"noNumberHere"`
	}

	bytes := []byte(`["value"]`)

	var out otherTag
	err := UnmarshalAsObject(bytes, &out)

	assert.Nil(t, err)
	assert.Equal(t, otherTag{}, out)
}
