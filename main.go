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
		return -1
	}

	err = parse.ParseAndProcess()
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
