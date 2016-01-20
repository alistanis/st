package sterrors

import "fmt"

var (
	Verbose = false
)

func MutuallyExclusiveParameters(p, p2 string) error {
	return fmt.Errorf("Mutually exclusive parameters provided: %s and %s", p, p2)
}

func Printf(s string, args ...interface{}) {
	if Verbose {
		fmt.Printf(s, args...)
	}
}
