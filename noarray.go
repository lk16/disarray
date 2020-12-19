package noarray

import (
	"encoding/json"
	"fmt"
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
	var rawSlice []interface{}
	var err error

	if err = json.Unmarshal(data, &rawSlice); err != nil {
		return err
	}

	x := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < x.NumField(); i++ {
		jsonTag := x.Type().Field(i).Tag.Get("json")
		jsonField := strings.SplitN(jsonTag, ",", 2)[0]

		var sliceIndex int
		if sliceIndex, err = strconv.Atoi(jsonField); err != nil {
			return err
		}

		if sliceIndex >= len(rawSlice) {
			// prevent index out of range
			continue
		}

		var fieldBytes []byte
		if fieldBytes, err = json.Marshal(rawSlice[sliceIndex]); err != nil {
			return fmt.Errorf("remarshalling at offset %d: %w", i, err)
		}

		if err = json.Unmarshal(fieldBytes, x.Field(i).Addr().Interface()); err != nil {
			return err
		}

	}
	return nil
}
