package noarray

import (
	"encoding/json"
	"errors"
)

var (
	Offsets = []string{
		"offset_zero",
		"offset_one",
		"offset_two",
		"offset_three",
		"offset_four",
		"offset_five",
		"offset_six",
		"offset_seven",
		"offset_eight",
		"offset_nine",
	}

	ErrArrayTooLong = errors.New("JSON array is longer than supported")
)

// UnmarshalAsObject unmarshals JSON array data into an object using keys
// named after their orignal offset in the array. This is useful when
// attempting to unmarshal JSON data that does not have an object as
// the root element.
//
// JSON data that is invalid or does not contain a JSON array will
// cause an error to be returned
//
// Arrays with more than 10 items are not supported.
//
// See tests for examples.
func UnmarshalAsObject(data []byte, v interface{}) error {
	var rawSlice []interface{}
	var err error

	if err = json.Unmarshal(data, &rawSlice); err != nil {
		return err
	}

	if len(rawSlice) > len(Offsets) {
		return ErrArrayTooLong
	}

	object := make(map[string]interface{}, len(rawSlice))

	for i := range rawSlice {
		object[Offsets[i]] = rawSlice[i]
	}

	var objectData []byte
	if objectData, err = json.Marshal(object); err != nil {
		return err
	}

	return json.Unmarshal(objectData, v)
}
