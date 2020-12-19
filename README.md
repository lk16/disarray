# No Array

This module exports one function with one task:

to of JSON messages that have an array as root element into a golang struct.

Especially those with different types as fields. This gets particularly hairy for complex JSON structures.

### Example

```go
import "github.com/lk16/noarray"

type foo struct {
	Zero string  `json:"offset_zero"`
	One  float64 `json:"offset_one"`
	Two  string  `json:"offset_two"`
}

func main() {
    input := []byte(`["foo",123.4,"bar"]`)

    var output foo
    err := noarray.UnmarshalAsObject(testCase.input, &output)

    if err != nil {
        panic(err.Error())
    }

    // prints: foo = foo{Zero: "foo", One: 123.4, Two: "bar"}
    fmt.Printf("foo = %+#v", foo)
}

```
