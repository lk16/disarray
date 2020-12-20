package disarray

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// UnmarshalAsObject unmarshals JSON array data into a golang struct.
//
// JSON data that is invalid or does not contain a JSON array will
// cause an error to be returned
//
// This function only affects the JSON array at JSON top level.
//
// Use tags like this to indicate which JSON array item should be
// unmarshalled for a field.
//  type foo struct {
//  	Zero string  `json:"0"`
//  	One  float64 `json:"1"`
//  	Two  string  `json:"2"`
//  }
//
// This function is most useful when called from UnmarshalJSON:
//  func (foo *foo) UnmarshalJSON(bytes []byte) error {
//  	return disarray.UnmarshalAsObject(bytes, foo)
//  }
//
func UnmarshalAsObject(data []byte, v interface{}) error {
	var rawMessages []json.RawMessage
	var err error

	if err = json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}

	x := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < x.NumField(); i++ {
		jsonTag, ok := x.Type().Field(i).Tag.Lookup("json")
		if !ok {
			continue
		}

		jsonField := strings.SplitN(jsonTag, ",", 2)[0]

		var sliceIndex int
		if sliceIndex, err = strconv.Atoi(jsonField); err != nil {
			return err
		}

		if sliceIndex >= len(rawMessages) {
			// prevent index out of range
			continue
		}

		unmarshalTarget := x.Field(i).Addr().Interface()
		if err = json.Unmarshal(rawMessages[i], unmarshalTarget); err != nil {
			return err
		}
	}
	return nil
}
