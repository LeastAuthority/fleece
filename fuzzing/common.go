package fuzzing

import "strings"

// Fuzz constants for go-fuzz to use when returning from the Fuzz func
var (
	FuzzInteresting = 1
	FuzzNormal      = 0
	FuzzDiscard     = -1
)

func FirstLine(str string) string {
	return strings.SplitN(str, "\n", 2)[0]
}