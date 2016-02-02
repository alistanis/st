package sterrors

import (
	"encoding/json"
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

type HttpError struct {
	Err  string `json:"error"`
	Code int    `json:"status_code"`
}

func FormatHTTPError(err error, code int) []byte {
	httpErr := &HttpError{Err: err.Error(), Code: code}
	// we bury this error because we know that the type passed to it will always be the right type
	data, _ := json.Marshal(httpErr)
	return data
}
