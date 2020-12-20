# Disarray

Unmarshal JSON arrays into a struct.

[Documentation](https://pkg.go.dev/github.com/lk16/disarray)

### Why?
* The Standard library doesn't support this
* Writing a custom `UnmarshalJSON()` for every array of some APIs object got tedious
* This seems reusable by many.

### Example

Full example is [here](example/main.go).

```go
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
    err := json.Unmarshal(bytes, &foo)
    // ...
}
```

### Discussion

Pros:
* This integrates nicely with existing code and nested JSON arrays.
* The tags used with indexes look clean

Cons:
* You need to implement a 3-line `UnmarshalJSON` for every JSON array model
* This could be supported by the standard library directly.
