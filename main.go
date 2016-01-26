package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alistanis/st/flags"
	"github.com/alistanis/st/net"
	"github.com/alistanis/st/parse"
	"github.com/gopherjs/gopherjs/js"
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
	if js.Global != nil {
		js.Global.Set("st", map[string]interface{}{
			"NewTagger": NewTagger,
		})
	} else {
		if flags.ServeStaticHttp {
			block := make(chan bool)
			var err error
			go func() {
				err = net.ServeStaticContent()
				if err != nil {
					fmt.Println(err)
					block <- true
				}
			}()
			<-block
			if err != nil {
				return -1
			}
		} else {
			err = parse.ParseAndProcessFiles(flag.Args())
			if err != nil {
				fmt.Println(err)
				return -1
			}
		}
	}
	return 0
}

type JSTagger struct{}

func NewTagger() *js.Object {
	return js.MakeWrapper(&JSTagger{})
}

func (t *JSTagger) Process(s string) string {
	data := []byte(s)
	if !strings.Contains(s, "package") {
		data = parse.Insert(data, []byte("package placeholder\n"), 0)
	}

	data, _ = parse.ProcessBytes(data, "placeholder.go")
	return string(data)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: st [flags] [path ...]\n")
	flag.PrintDefaults()
	exit(-2)
}

func main() {
	exit(run())
}
