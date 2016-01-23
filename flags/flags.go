package flags

import (
	"errors"
	"flag"
	"strings"

	"github.com/alistanis/st/parse"
	"github.com/alistanis/st/sterrors"
)

var (
	Case = parse.DefaultCase
	Tag  = parse.DefaultTag

	Append    bool
	Overwrite bool
	c         bool
	s         bool
	Verbose   bool
	Write     bool

	IgnoredFieldsString  string
	IgnoredStructsString string

	AppendMode = parse.SkipExisting
	TagMode    = parse.TagAll
)

const (
	Camel = "camel"
	Snake = "snake"
)

func StringVars() {
	flag.StringVar(&Tag, "t", "json", "The struct tag to use when tagging. Example: -t=json ")
	flag.StringVar(&IgnoredFieldsString, "i", "", "A comma separated list of fields to ignore. Will use the format json:\"-\".")
	flag.StringVar(&IgnoredStructsString, "is", "", "A comma separated list of structs to ignore. Will not tag any fields in the struct.")

}

func BoolVars() {
	flag.BoolVar(&c, "c", false, "Sets the struct tag to camel case.")
	flag.BoolVar(&s, "s", false, "Sets the struct tag to snake case.")
	flag.BoolVar(&Append, "a", false, "Sets mode to append mode. Will append to existing tags. Default behavior skips existing tags.")
	flag.BoolVar(&Verbose, "v", false, "Sets mode to verbose.")
	flag.BoolVar(&Write, "w", false, "Sets mode to write to source file.")
	flag.BoolVar(&Overwrite, "o", false, "Sets mode to overwrite mode. Will overwrite existing tags (completely). Default behavior skips existing tags.")
}

func SetVars() {
	StringVars()
	BoolVars()
}

func ParseFlags() error {
	SetVars()
	flag.Parse()
	return verify()
}

func verify() error {
	if flag.NArg() < 1 {
		return errors.New("No path was provided.")
	}

	if c && s {
		return sterrors.MutuallyExclusiveParameters("c", "s")
	}

	if Overwrite && Append {
		return sterrors.MutuallyExclusiveParameters("o", "a")
	}

	if c {
		Case = Camel
	}

	if s {
		Case = Snake
	}

	if Overwrite {
		AppendMode = parse.Overwrite
	}

	if Append {
		AppendMode = parse.Append
	}

	sterrors.Verbose = Verbose

	if IgnoredFieldsString != "" {
		parse.IgnoredFields = strings.Split(IgnoredFieldsString, ",")
	}

	if IgnoredStructsString != "" {
		parse.IgnoredStructs = strings.Split(IgnoredStructsString, ",")
	}
	return nil
}
