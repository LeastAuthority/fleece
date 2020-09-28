//+build gofuzz

package example

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	fleece "github.com/leastauthority/fleece/fuzzing"
)

var (
	limit       int
	fleeceDir   string
	skipPattern string

	env *fleece.Env
)

func init() {
	flag.IntVar(&limit, "crash-limit", 1000, "number of crashing inputs to test before stopping")
	flag.StringVar(&fleeceDir, "fleece-dir", "fleece", "path to fleece dir relative to repo/module root")
	flag.StringVar(&skipPattern, "skip", "", "if provided, crashers with recorded outputs which match the pattern will be skipped")
}

func TestMain(m *testing.M) {
	flag.Parse()
	env = fleece.NewEnv(fleeceDir)

	os.Exit(m.Run())
}

func TestFuzzPanickyFunc(t *testing.T) {
	_, panics, _ := fleece.MustNewCrasherIterator(env, FuzzPanickyFunc).
		TestFailingLimit(t, crashLimit, fleece.SkipFilter(skipPattern))

	require.Zero(t, panics)
}
