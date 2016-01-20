package flags

import (
	"errors"
	"flag"

	"github.com/alistanis/st/sterrors"
)

var (
	Case string
	Tag  string

	Append  bool
	c       bool
	s       bool
	Verbose bool
	Write   bool
)

const (
	Camel = "camel"
	Snake = "snake"
)

func StringVars() {
	flag.StringVar(&Tag, "t", "json", "The struct tag to use when tagging. Example: `json:\"var_name\"`")
}

func BoolVars() {
	flag.BoolVar(&c, "c", false, "Sets the struct tag to camel case")
	flag.BoolVar(&s, "s", false, "Sets the struct tag to snake case")
	flag.BoolVar(&Append, "a", false, "Sets mode to append mode. The default is to overwrite existing tags.")
	flag.BoolVar(&Verbose, "v", false, "Sets mode to verbose. (prints extra information)")
	flag.BoolVar(&Write, "w", false, "Sets mode to write to source file.")
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
		return errors.New("No path was provided. The -path parameter is required.")
	}

	if c && s {
		return sterrors.MutuallyExclusiveParameters("c", "s")
	}

	if c {
		Case = Camel
	}

	if s {
		Case = Snake
	}
	sterrors.Verbose = Verbose
	return nil
}
