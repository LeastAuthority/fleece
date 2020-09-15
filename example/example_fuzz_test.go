//+build gofuzz

package example

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/leastauthority/lafuzz/fuzzing"
)

func TestFuzzExample(t *testing.T) {
	_, panics, _ := fuzzing.
		MustNewCrasherIteratorFor(FuzzExample).
		TestFailingLimit(t, 1000)

	require.Zero(t, panics)
}
