package fuzzing

import (
	"fmt"
	"runtime/debug"
)

// Crasher represents a go-fuzz "crasher" (an input that crashed the respective
//	fuzz function), its corresponding output (panic message), and name (input hash).
type Crasher struct {
	Name     string
	Input    []byte
	//Output   string
	FuzzFunc Func
}

type RecoverCallback func(panicMsg string)

// Recover is intended to be deferred. It calls the recover callback with the
//	string representation of the recovered value in the event of a panic.
func (crasher *Crasher) Recover(recoverCb RecoverCallback) {
	if r := recover(); r != nil {
		recoverCb(fmt.Sprintf("%s\n%s", r, string(debug.Stack())))
	}
}

// Test calls the crashers fuzz function its input and recovers from panics
//	with the passed recover callback.
func (crasher *Crasher) Test(recoverCb RecoverCallback) {
	defer crasher.Recover(recoverCb)
	crasher.FuzzFunc(crasher.Input)
}

