package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alistanis/st/flags"
	"github.com/alistanis/st/parse"
)

func run() int {
	flag.Usage = usage
	err := flags.ParseFlags()
	if err != nil {
		fmt.Println(err)
		usage()
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
	os.Exit(2)
}

func main() {
	os.Exit(run())
}
