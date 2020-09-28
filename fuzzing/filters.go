package fuzzing

import (
	"fmt"
	"strings"
)

const TimedOutPattern = "program hanged"

var SkipTimedOut = SkipFilter(TimedOutPattern)

func SkipFilter(pattern string) IterFilter {
	if pattern == "" {
		return all
	}

	return func(next *Crasher) bool {
		firstLine := strings.SplitN(next.Output, "\n", 2)[0]
		skip := strings.Contains(firstLine, pattern)
		if skip {
			fmt.Printf("skipping %s, output matches pattern %q\n", next.Name, pattern)
		}
		return !skip
	}
}

func all(_ *Crasher) bool {
	return true
}
