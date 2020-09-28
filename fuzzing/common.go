package fuzzing

// Fuzz constants for go-fuzz to use when returning from the Fuzz func
var (
	FuzzInteresting = 1
	FuzzNormal      = 0
	FuzzDiscard     = -1
)
