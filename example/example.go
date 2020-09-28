package example

import (
	"fmt"
)

func PanickyFunc(input []byte) {
	size := len(input)
	if 5 < size && size < 13 {
		panic(fmt.Errorf("input size: %d", size))
	}
}
