package sterrors

import (
	"errors"
	"fmt"
)

var (
	// Verbose determines whether or not Printf will print anything
	Verbose = false
	// ErrNoPathsGiven is returned when no paths to any .go files were provided at the command line
	ErrNoPathsGiven = errors.New("No paths to any .go files were provided.")
)

// ErrMutuallyExclusiveParameters takes two inputs and returns a canned error response
func ErrMutuallyExclusiveParameters(p, p2 string) error {
	return fmt.Errorf("Mutually exclusive parameters provided: %s and %s", p, p2)
}

// Printf prints a string depending on verbosity... should be in a debug package?
func Printf(s string, args ...interface{}) {
	if Verbose {
		fmt.Printf(s, args...)
	}
}
