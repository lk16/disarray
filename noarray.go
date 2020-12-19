package noarray

import (
	"encoding/json"
	"fmt"
)

// UnmarshalAsObject unmarshals JSON array data into a golang struct.
// This is useful when attempting to unmarshal a JSON message
// that has a JSON array as root.
//
// JSON data that is invalid or does not contain a JSON array will
// cause an error to be returned
//
// See tests and repo readme for examples.
func UnmarshalAsObject(data []byte, v interface{}) error {
	var rawSlice []interface{}
	var err error

	if err = json.Unmarshal(data, &rawSlice); err != nil {
		return err
	}

	object := make(map[string]interface{}, len(rawSlice))

	for i := range rawSlice {
		object[fmt.Sprintf("%d", i)] = rawSlice[i]
	}

	var objectData []byte
	if objectData, err = json.Marshal(object); err != nil {
		return err
	}

	return json.Unmarshal(objectData, &v)
}
