ST - Struct Tagger for Go
---
st is a command line utility for tagging structs in your Go code.

[![Build Status](https://travis-ci.org/alistanis/st.svg?branch=master)](https://travis-ci.org/alistanis/st) ![Report Card](http://goreportcard.com/badge/alistanis/st) [![Gitter](https://badges.gitter.im/alistanis/st.svg)](https://gitter.im/alistanis/st?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
---

Get it: 
```go get github.com/alistanis/st```

If you want to run the tests, you'll need the goconvey package:
```go get github.com/smartystreets/goconvey```

Roadmap
---
```
1. [-] Command Line Support
	* [x] Write tags to buffer(default behavior) and to file (with -w)
	* [x] Supports multiple operation modes
		* [x] Append
		* [x] Overwrite
		* [x] Skip existing tags
		* [x] Field Exclusion
		* [x] Struct Exclusion
		* [ ] Explicit Struct/field inclusion
		* [ ] Go Generate support 
2. [ ] Web Application
	* [x] Basic static site handler (not in master)
	* [ ] Side by side input/output
	* [ ] Make it pretty  
3. [x] Tests/Build/Deploy
	* [x] Main Package tests
	* [x] Parse Package tests
	* [x] Flags package tests
	* [x] Travis Integration
		* [x] Run tests automatically 
		* [x] Build notifications
4. [x] Miscellaneous
	* [x] Gitter Integration
	* [x] Slack Channel
```

Usage
---
>```st -h or st --help```

```
usage: st [flags] [path ...]
  -a	Sets mode to append mode. Will append to existing tags. Default behavior skips existing tags.
  -append
    	Sets mode to append mode. Will append to existing tags. Default behavior skips existing tags.
  -c	Sets the struct tag to camel case.
  -camel
    	Sets the struct tag to camel case
  -i string
    	A comma separated list of fields to ignore. Will use the format json:"-".
  -ignored-fields string
    	A comma separated list of fields to ignore. Will use the format json:"-".
  -ignored-structs string
    	A comma separated list of structs to ignore. Will not tag any fields in the struct.
  -is string
    	A comma separated list of structs to ignore. Will not tag any fields in the struct.
  -o	Sets mode to overwrite mode. Will overwrite existing tags (completely). Default behavior skips existing tags.
  -overwrite
    	Sets mode to overwrite mode. Will overwrite existing tags (completely). Default behavior skips existing tags.
  -s	Sets the struct tag to snake case.
  -snake
    	Sets the struct tag to snake case.
  -t string
    	The struct tag to use when tagging. Example: -t=json  (default "json")
  -tag-name string
    	The struct tag to use when tagging. Example: --tag-name=json  (default "json")
  -v	Sets mode to verbose.
  -verbose
    	Sets mode to verbose.
  -w	Sets mode to write to source file. The default is a dry run that prints the results to stdout.
  -write
    	Sets mode to write to source file. The default is a dry run that prints the results to stdout.
```

Defaults
---
>
* ST will not write to your source file unless you provide the **-w** or **-write** flags. Its default behavior prints the result to *STDOUT*
* The default tag that ST uses is **json**
* The default tagging mode is to *Skip Existing Tags* - you can change this behavior by providing one of the *Append* flags, **-a** or **-append**, or by using one of the *Overwrite* flags, **-o** or **-overwrite**
* The default tagging case is *Snake Case* - this can be changed by providing either *Camel Case* flag, **-c** or **-camel**  
>
>Overwrite mode will completely overwrite an existing tag. Append mode is a little trickier. If an existing tag is there for the
tag that you have specified, let's use json as our example, it will leave that tag alone. If you specify a different tag, like msgpack,
it will append to the existing tag with the msgpack key/value.

Overwrite Examples 
---
>```st --overwrite --tag-name=msgpack $GOFILE```

```go
type Test struct { F field `json:"f"`}
    becomes
type Test struct { F field `msgpack:"f"`}
```
>```st --overwrite --tag-name=json --case=camel $GOFILE```

```go
type Test struct { F field `json:"f"`}
    becomes
type Test struct { F field `json:"F"`}
```

Append Examples
---
>```st --append --case=camel --tag-name=json $GOFILE```

```go
type Test struct { F field `json:"f"`}
    becomes (the tag is left alone because it is already there)
type Test struct { F field `json:"f"`}
```
>```st --append --tag-name=msgpack $GOFILE```

```go
type Test struct { F field `json:"f"`}
    becomes
type Test struct { F field `msgpack:"f" json:"f"`}
```

Contributing & Contact
---
If you would like to contribute, don't be shy! Fork the project, write tests for any new code and ensure that you don't break existing
functionality by running the current tests. If you're looking to submit changes upstream, it would be a good idea to
discuss it with me first. I'm available via [email](ccooper@sessionm.com), through Github, on the
[Gophers Slack ST Channel](https://blog.gopheracademy.com/gophers-slack-community/), or on [Gitter.im](https://gitter.im/alistanis/st).

If you do submit a pull request, I will review it and I will merge it if it's in line with my vision for the project.


Further examples
---

>Contents of etc.go before running any of the following

```go
package etc

type TestStruct struct {
	Int             int
	Int64           int64
	IntSlice        []int
	...
}

type TestStructWithTagsSnake struct {
	Int             int               `json:"int"`
	Int64           int64             `json:"int_64"`
	IntSlice        []int             `json:"int_slice"`
	...
}

type TestStructWithTagsCamel struct {
	Int             int               `json:"Int"`
	Int64           int64             `json:"Int64"`
	...
}

type TestUnexportedField struct {
	unexportedField int
	ExportedField   int
}

```

> Append to existing tags with the tag msgpack (use -w flag to write to original source file) 
```
st -s -a -v -t=msgpack $GOPATH/src/github.com/alistanis/st/etc/etc.go
```

```go 
package etc

type TestStruct struct {
	Int             int               `msgpack:"int"`
	Int64           int64             `msgpack:"int_64"`
	IntSlice        []int             `msgpack:"int_slice"`
	...
}

type TestStructWithTagsSnake struct {
	Int             int               `msgpack:"int" json:"int"`
	Int64           int64             `msgpack:"int_64" json:"int_64"`
	IntSlice        []int             `msgpack:"int_slice" json:"int_slice"`
	...
}

type TestStructWithTagsCamel struct {
	Int             int               `msgpack:"int" json:"Int"`
	Int64           int64             `msgpack:"int_64" json:"Int64"`
	IntSlice        []int             `msgpack:"int_slice" json:"IntSlice"`
	...
}

type TestUnexportedField struct {
	unexportedField int
	ExportedField   int `msgpack:"exported_field"`
}
```
>Ignore a specific field (-i) and ignore a specific struct (-is)
```
st -s -a -v -i=ExportedField -is=TestStructWithTagsCamel -t=msgpack $GOPATH/src/github.com/alistanis/st/etc/etc.go
```

```go
package etc

type TestStruct struct {
	Int             int               `msgpack:"int"`
	Int64           int64             `msgpack:"int_64"`
	IntSlice        []int             `msgpack:"int_slice"`
	...
}

type TestStructWithTagsSnake struct {
	Int             int               `msgpack:"int" json:"int"`
	Int64           int64             `msgpack:"int_64" json:"int_64"`
	IntSlice        []int             `msgpack:"int_slice" json:"int_slice"`
	...
}

type TestStructWithTagsCamel struct {
	Int             int               `json:"Int"`
	Int64           int64             `json:"Int64"`
	IntSlice        []int             `json:"IntSlice"`
	...
}

type TestUnexportedField struct {
	unexportedField int
	ExportedField   int `msgpack:"-"`
}
```
