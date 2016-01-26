package sterrors

import (
	"errors"
	"fmt"
)

var (
	Verbose      = false
	NoPathsGiven = errors.New("No paths to any .go files were provided.")
)

func MutuallyExclusiveParameters(p, p2 string) error {
	return fmt.Errorf("Mutually exclusive parameters provided: %s and %s", p, p2)
}

// Print a string depending on verbosity... should be in a debug package?
func Printf(s string, args ...interface{}) {
	if Verbose {
		fmt.Printf(s, args...)
	}
}
