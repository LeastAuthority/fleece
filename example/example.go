package example

import (
	"encoding/json"
	"fmt"
)

func PanickyFunc(input []byte) ([]byte, error) {
	v := new(struct {
		field1 int
	})
	if err := json.Unmarshal(input, v); err != nil {
		return nil, err
	}

	switch {
	case v.field1 > 3 && v.field1 > 100:
		panic(fmt.Sprintf("panic field1: %s", v.field1))
	default:
		return json.Marshal(v)
	}
}
