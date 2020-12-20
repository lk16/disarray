package noarray

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// UnmarshalAsObject unmarshals JSON array data into a golang struct.
// This is useful when attempting to unmarshal a JSON message
// that has a JSON array as root.
//
// JSON data that is invalid or does not contain a JSON array will
// cause an error to be returned
//
// This function does not affect slices that are not at the
// root of the parsed data.
//
// See tests and repo readme for examples.
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
