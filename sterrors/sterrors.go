package sterrors

import "fmt"

var (
	Verbose = false
)

func MutuallyExclusiveParameters(p, p2 string) error {
	return fmt.Errorf("Mutually exclusive parameters provided: %s and %s", p, p2)
}

// Print a string depending on verbosity
func Printf(s string, args ...interface{}) {
	if Verbose {
		fmt.Printf(s, args...)
	}
}
