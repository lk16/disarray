package main

import (
	"encoding/json"
	"fmt"

	"github.com/lk16/disarray"
)

type foo struct {
	Zero string  `json:"0"`
	One  float64 `json:"1"`
	Two  string  `json:"2"`
}

func (foo *foo) UnmarshalJSON(bytes []byte) error {
	return disarray.UnmarshalAsObject(bytes, foo)
}

func main() {
	bytes := []byte(`["foo",123.4,"bar"]`)

	var foo foo
	if err := json.Unmarshal(bytes, &foo); err != nil {
		panic(err.Error())
	}

	fmt.Printf("foo = %+#v\n", foo)
}
