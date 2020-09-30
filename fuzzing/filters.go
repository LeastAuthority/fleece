package fuzzing

import (
	"fmt"
	"strings"
)

const TimedOutPattern = "program hanged"
const OutOfMemoryPattern = "out of memory"

var (
	SkipTimedOut    = SkipFilter(TimedOutPattern, "", false)
	SkipOutOfMemory = SkipFilter(OutOfMemoryPattern, "", false)
)

type IterFilter func(next *Crasher) bool
type IterFilters []IterFilter

func SkipFilter(patternStr string, delimiter string, verbose bool) IterFilter {
	filtered := 0
	if patternStr == "" {
		return all
	}

	var patterns []string
	if delimiter == "" {
		patterns = []string{patternStr}
	} else {
		patterns = strings.Split(patternStr, delimiter)
	}

	return func(next *Crasher) bool {
		for _, pattern := range patterns {
			skip := strings.Contains(FirstLine(next.Output), pattern)
			if skip {
				filtered++
				if filtered%1000 == 0 || verbose {
					fmt.Printf("skip %d: %s, output matches pattern %q\n", filtered, next.Name, pattern)
				}
				return false
			}
		}
		return true
	}
}

func (filters IterFilters) Allows(next *Crasher) bool {
	for _, filter := range filters {
		if !filter(next) {
			return false
		}
	}
	return true
}

func all(_ *Crasher) bool {
	return true
}
