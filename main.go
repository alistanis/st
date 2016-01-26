package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alistanis/st/flags"
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
	err := flags.ParseFlags()
	if err != nil {
		fmt.Println(err)
		usage()
		return -1
	}

	options := &parse.Options{
		Tag:        flags.Tag,
		Case:       flags.Case,
		AppendMode: flags.AppendMode,
		TagMode:    flags.TagMode,
		// this is confusing, I'll fix it later when changing documentation/flags behavior
		DryRun:  !flags.Write,
		Verbose: flags.Verbose}
	parse.SetOptions(options)
	err = parse.ParseAndProcessFiles(flag.Args())
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return 0
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: st [flags] [path ...]\n")
	flag.PrintDefaults()
	exit(-2)
}

func main() {
	exit(run())
}
