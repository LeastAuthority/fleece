//+build gofuzz

package example

import (
	"github.com/leastauthority/fleece/fuzzing"
)

func FuzzPanickyFunc(data []byte) int {
	PanickyFunc(data)

	return fuzzing.FuzzNormal
}
