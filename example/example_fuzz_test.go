//+build gofuzz

package example

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/leastauthority/lafuzz/fuzzing"
)

func TestFuzzExample(t *testing.T) {
	_, panics, _ := fuzzing.
		MustNewCrasherIterator(FuzzExample).
		TestFailingLimit(t, 1000)

	require.Zero(t, panics)
}
