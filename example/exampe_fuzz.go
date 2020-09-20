//+build gofuzz

package example

import (
	"bytes"

	"github.com/leastauthority/fleece/fuzzing"
)

func FuzzBuggyFunc(data []byte) int {
	result, err := PanickyFunc(data)
	if err != nil {
		return fuzzing.FuzzNormal
	}

	if !bytes.Equal(result, data) {
		panic("input and result aren't equal!")
	}

	return 1
}
