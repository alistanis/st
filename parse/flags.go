package parse

import (
	"flag"
	"os"
	"strings"

	"github.com/alistanis/st/sterrors"
)

var (
	// Case determines the case to use when tagging structs - either Camel or Snake
	Case = DefaultCase
	// Tag determines the tag to use when tagging structs - default is json
	Tag = DefaultTag
	// FlagAppend is true if -a or -append are provided as command line flags - appends to tags instead of overwriting or skipping entirely
	FlagAppend bool
	// FlagOverwrite is true if -o or -overwrite are provided as command line flags - overwrites existing tags
	FlagOverwrite bool
	c             bool
	s             bool
	// Verbose sets the default for how much information is printed to standard out
	Verbose bool
	// Write is true if -w or -write are provided as command line flags - this will write to the original source file
	Write bool
	// IgnoredFieldsString is a comma separated list of ignored fields provided as a command line flag
	IgnoredFieldsString string
	// IgnoredStructsString is a comma separated list of ignored structs provided as a command line flag
	IgnoredStructsString string
	// AppendMode is the mode that ST will operate in. Default is to skip existing tags, can be set to Append or Overwrite
	AppendMode = SkipExisting
	// TagMode is the mode that ST operates on when tagging. Default is to tag all structs/fields.
	TagMode = TagAll
	// GoFile is the name of the GoFile as given by go generate to os.Environ ($GOFILE)
	GoFile string
)

// stringVars sets up all string command line variable bindings
func stringVars() {
	flag.StringVar(&Tag, "t", "json", "The struct tag to use when tagging. Example: -t=json ")
	flag.StringVar(&Tag, "tag-name", "json", "The struct tag to use when tagging. Example: --tag-name=json ")
	flag.StringVar(&IgnoredFieldsString, "i", "", "A comma separated list of fields to ignore. Will use the format json:\"-\".")
	flag.StringVar(&IgnoredFieldsString, "ignored-fields", "", "A comma separated list of fields to ignore. Will use the format json:\"-\".")
	flag.StringVar(&IgnoredStructsString, "is", "", "A comma separated list of structs to ignore. Will not tag any fields in the struct.")
	flag.StringVar(&IgnoredStructsString, "ignored-structs", "", "A comma separated list of structs to ignore. Will not tag any fields in the struct.")
}

// boolVars sets up all boolean command line variable bindings
func boolVars() {
	flag.BoolVar(&c, "c", false, "Sets the struct tag to camel case.")
	flag.BoolVar(&c, "camel", false, "Sets the struct tag to camel case")
	flag.BoolVar(&s, "s", false, "Sets the struct tag to snake case.")
	flag.BoolVar(&s, "snake", false, "Sets the struct tag to snake case.")
	flag.BoolVar(&FlagAppend, "a", false, "Sets mode to append mode. Will append to existing tags. Default behavior skips existing tags.")
	flag.BoolVar(&FlagAppend, "append", false, "Sets mode to append mode. Will append to existing tags. Default behavior skips existing tags.")
	flag.BoolVar(&Verbose, "v", false, "Sets mode to verbose.")
	flag.BoolVar(&Verbose, "verbose", false, "Sets mode to verbose.")
	flag.BoolVar(&Write, "w", false, "Sets mode to write to source file. The default is a dry run that prints the results to stdout.")
	flag.BoolVar(&Write, "write", false, "Sets mode to write to source file. The default is a dry run that prints the results to stdout.")
	flag.BoolVar(&FlagOverwrite, "o", false, "Sets mode to overwrite mode. Will overwrite existing tags (completely). Default behavior skips existing tags.")
	flag.BoolVar(&FlagOverwrite, "overwrite", false, "Sets mode to overwrite mode. Will overwrite existing tags (completely). Default behavior skips existing tags.")
}

// SetVars sets up all command line variable bindings
func SetVars() {
	stringVars()
	boolVars()
	GoFile = os.Getenv("GOFILE")
}

// Flags sets up command line bindings, calls flagParse(), and calls verify() to check command line flags
func Flags() error {
	SetVars()
	flag.Parse()
	return verify()
}

func verify() error {

	// If GoFile is set, we know that we're being run by go generate, so we append the file name as our last argument
	// and we cheat so that we get the desired behavior
	if GoFile != "" {
		SetArgs(append(flag.Args(), GoFile))
		flag.Parse()
	}

	if flag.NArg() < 1 {
		return sterrors.ErrNoPathsGiven
	}

	if c && s {
		return sterrors.ErrMutuallyExclusiveParameters("c", "s")
	}

	if FlagOverwrite && FlagAppend {
		return sterrors.ErrMutuallyExclusiveParameters("o", "a")
	}

	if c {
		Case = Camel
	}

	if s {
		Case = Snake
	}

	if FlagOverwrite {
		AppendMode = Overwrite
	}

	if FlagAppend {
		AppendMode = Append
	}

	sterrors.Verbose = Verbose

	if IgnoredFieldsString != "" {
		IgnoredFields = strings.Split(IgnoredFieldsString, ",")
	}

	if IgnoredStructsString != "" {
		IgnoredStructs = strings.Split(IgnoredStructsString, ",")
	}
	return nil
}

// ResetFlags is a near copy of the flag.ResetForTesting(usage func()) function.
func ResetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

// SetArgs clears flags and sets os.Args to os.Args[0] (program name) and then to the list of whatever parameters are given after
func SetArgs(s []string) {
	ResetFlags()
	os.Args = []string{os.Args[0]}
	os.Args = append(os.Args, s...)
}
