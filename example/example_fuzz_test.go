//+build gofuzz

package example

import (
	"flag"
	"github.com/leastauthority/fleece/fuzzing"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	crashLimit int
	fleeceDir  string
	env        *fuzzing.Env
)

func init() {
	flag.IntVar(&crashLimit, "crash-limit", 1000, "number of crashing inputs to test before stopping")
	flag.StringVar(&fleeceDir, "fleece-dir", "./fleece", "path to the root of fuzz function workdirs")
}

func TestMain(m *testing.M) {
	flag.Parse()
	env = fuzzing.NewEnv(fleeceDir)

	os.Exit(m.Run())
}

func TestFuzzPanickyFunc(t *testing.T) {
	_, panics, _ := fuzzing.MustNewCrasherIterator(env, FuzzPanickyFunc).
		TestFailingLimit(t, env, crashLimit)

	require.Zero(t, panics)
}
