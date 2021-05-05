//+build gofuzz

package example

import (
	"github.com/leastauthority/fleece/fuzzing"
)

func FuzzBuggyFunc(data []byte) int {
	BuggyFunc(data)

	return fuzzing.FuzzNormal
}
