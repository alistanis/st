package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alistanis/st/parse"
)

var (
	exitFunction = defaultExitFunc
)

func defaultExitFunc(code int) {
	os.Exit(code)
}

func exit(code int) {
	if exitFunction != nil {
		exitFunction(code)
	} else {
		defaultExitFunc(code)
	}
}

func run() int {
	flag.Usage = usage
	err := parse.Flags()
	if err != nil {
		fmt.Println(err)
		usage()
		return -1
	}

	options := &parse.Options{
		Tag:        parse.Tag,
		Case:       parse.Case,
		AppendMode: parse.AppendMode,
		TagMode:    parse.TagMode,
		// this is confusing, I'll fix it later when changing documentation/flags behavior
		DryRun:  !parse.Write,
		Verbose: parse.Verbose}
	parse.SetOptions(options)
	err = parse.AndProcessFiles(flag.Args())
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return 0
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: st [flags] [path ...]")
	flag.PrintDefaults()
	exit(-2)
}

func main() {
	exit(run())
}
