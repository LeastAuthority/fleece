package fuzzing

import (
	"fmt"
	"strings"
)

const TimedOutPattern = "program hanged"

func SkipFilter(pattern string) IterFilter {
	return func(next *Crasher) bool {
		firstLine := strings.SplitN(next.Output, "\n", 2)[0]
		skip := strings.Contains(firstLine, pattern)
		if skip {
			fmt.Printf("skipping %s, output matches pattern %q\n", next.Name, pattern)
		}
		return !skip
	}
}

// TODO: better UX
func SkipTimedOut() IterFilter {
	skip := SkipFilter(TimedOutPattern)
	return func(next *Crasher) bool {
		return skip(next)
	}
}
