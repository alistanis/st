# st
ST - Struct Tagger for Go

st is a command line utility for tagging structs in your Go code.

---

Get it: 
```go get github.com/alistanis/st```

If you want to run the tests, you'll need the goconvey package:
```go get github.com/smartystreets/goconvey```

```
st -h or st --help
```

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

#Notes:

Defaults are set to what I consider to be reasonable. This means that ST won't write to your source file unless you
provide the -w or --write flags, it will simply print what the result will be to your buffer. The default tag it will use
if no tag is specified is json, and the default case it uses is snake case. ST, by default, will skip any struct fields
that already have tags. That behavior can be overridden by using the overwrite mode (-o or --overwrite), or append mode (-a or --append).

Overwrite mode will completely overwrite an existing tag. Append mode is a little trickier. If an existing tag is there for the
tag that you have specified, let's use json as our example, it will leave that tag alone. If you specify a different tag, like msgpack,
it will append to the existing tag with the msgpack key/value.


#Overwrite Examples: 
```
st --overwrite --tag-name=msgpack $GOFILE

type Test struct { F field `json:"f"`}
```
becomes
```
type Test struct { F field `msgpack:"f"`}
```

```
st --overwrite --tag-name=json --case=camel $GOFILE
type Test struct { F field `json:"f"`}
```
becomes
```
type Test struct { F field `json:"F"`}
```

#Append Examples:
```
st --append --case=camel --tag-name=json $GOFILE
type Test struct { F field `json:"f"`}
```
becomes (the tag is left alone because it is already there)
```
type Test struct { F field `json:"f"`}
```

```
st --append --tag-name=msgpack $GOFILE
type Test struct { F field `json:"f"`}
```
becomes
```
type Test struct { F field `msgpack:"f" json:"f"`}
```


Further examples

Contents of etc.go before running
```go
package etc

type TestStruct struct {
	Int             int
	Int64           int64
	IntSlice        []int
	Int64Slice      []int64
	String          string
	StringSlice     []string
	Float           float64
	FloatSlice      []float64
	UIntPointer     uintptr
	Rune            rune
	RuneSlice       []rune
	Byte            byte
	ByteSlice       []byte
	MapStringString map[string]string
	MapStringInt    map[string]int
	MapIntString    map[int]string
}

type TestStructWithTagsSnake struct {
	Int             int               `json:"int"`
	Int64           int64             `json:"int_64"`
	IntSlice        []int             `json:"int_slice"`
	Int64Slice      []int64           `json:"int_64_slice"`
	String          string            `json:"string"`
	StringSlice     []string          `json:"string_slice"`
	Float           float64           `json:"float"`
	FloatSlice      []float64         `json:"float_slice"`
	UIntPointer     uintptr           `json:"u_int_pointer"`
	Rune            rune              `json:"rune"`
	RuneSlice       []rune            `json:"rune_slice"`
	Byte            byte              `json:"byte"`
	ByteSlice       []byte            `json:"byte_slice"`
	MapStringString map[string]string `json:"map_string_string"`
	MapStringInt    map[string]int    `json:"map_string_int"`
	MapIntString    map[int]string    `json:"map_int_string"`
}

type TestStructWithTagsCamel struct {
	Int             int               `json:"Int"`
	Int64           int64             `json:"Int64"`
	IntSlice        []int             `json:"IntSlice"`
	Int64Slice      []int64           `json:"Int64Slice"`
	String          string            `json:"String"`
	StringSlice     []string          `json:"StringSlice"`
	Float           float64           `json:"Float"`
	FloatSlice      []float64         `json:"FloatSlice"`
	UIntPointer     uintptr           `json:"UIntPointer"`
	Rune            rune              `json:"Rune"`
	RuneSlice       []rune            `json:"RuneSlice"`
	Byte            byte              `json:"Byte"`
	ByteSlice       []byte            `json:"ByteSlice"`
	MapStringString map[string]string `json:"MapStringString"`
	MapStringInt    map[string]int    `json:"MapStringInt"`
	MapIntString    map[int]string    `json:"MapIntString"`
}

type TestUnexportedField struct {
	unexportedField int
	ExportedField   int
}

```

- Append to existing tags with the tag msgpack (use -w flag to write to original source file)

```
st -s -a -v -t=msgpack $GOPATH/src/github.com/alistanis/st/etc/etc.go
```

```go
package etc

type TestStruct struct {
	Int             int               `msgpack:"int"`
	Int64           int64             `msgpack:"int_64"`
	IntSlice        []int             `msgpack:"int_slice"`
	Int64Slice      []int64           `msgpack:"int_64_slice"`
	String          string            `msgpack:"string"`
	StringSlice     []string          `msgpack:"string_slice"`
	Float           float64           `msgpack:"float"`
	FloatSlice      []float64         `msgpack:"float_slice"`
	UIntPointer     uintptr           `msgpack:"u_int_pointer"`
	Rune            rune              `msgpack:"rune"`
	RuneSlice       []rune            `msgpack:"rune_slice"`
	Byte            byte              `msgpack:"byte"`
	ByteSlice       []byte            `msgpack:"byte_slice"`
	MapStringString map[string]string `msgpack:"map_string_string"`
	MapStringInt    map[string]int    `msgpack:"map_string_int"`
	MapIntString    map[int]string    `msgpack:"map_int_string"`
}

type TestStructWithTagsSnake struct {
	Int             int               `msgpack:"int" json:"int"`
	Int64           int64             `msgpack:"int_64" json:"int_64"`
	IntSlice        []int             `msgpack:"int_slice" json:"int_slice"`
	Int64Slice      []int64           `msgpack:"int_64_slice" json:"int_64_slice"`
	String          string            `msgpack:"string" json:"string"`
	StringSlice     []string          `msgpack:"string_slice" json:"string_slice"`
	Float           float64           `msgpack:"float" json:"float"`
	FloatSlice      []float64         `msgpack:"float_slice" json:"float_slice"`
	UIntPointer     uintptr           `msgpack:"u_int_pointer" json:"u_int_pointer"`
	Rune            rune              `msgpack:"rune" json:"rune"`
	RuneSlice       []rune            `msgpack:"rune_slice" json:"rune_slice"`
	Byte            byte              `msgpack:"byte" json:"byte"`
	ByteSlice       []byte            `msgpack:"byte_slice" json:"byte_slice"`
	MapStringString map[string]string `msgpack:"map_string_string" json:"map_string_string"`
	MapStringInt    map[string]int    `msgpack:"map_string_int" json:"map_string_int"`
	MapIntString    map[int]string    `msgpack:"map_int_string" json:"map_int_string"`
}

type TestStructWithTagsCamel struct {
	Int             int               `msgpack:"int" json:"Int"`
	Int64           int64             `msgpack:"int_64" json:"Int64"`
	IntSlice        []int             `msgpack:"int_slice" json:"IntSlice"`
	Int64Slice      []int64           `msgpack:"int_64_slice" json:"Int64Slice"`
	String          string            `msgpack:"string" json:"String"`
	StringSlice     []string          `msgpack:"string_slice" json:"StringSlice"`
	Float           float64           `msgpack:"float" json:"Float"`
	FloatSlice      []float64         `msgpack:"float_slice" json:"FloatSlice"`
	UIntPointer     uintptr           `msgpack:"u_int_pointer" json:"UIntPointer"`
	Rune            rune              `msgpack:"rune" json:"Rune"`
	RuneSlice       []rune            `msgpack:"rune_slice" json:"RuneSlice"`
	Byte            byte              `msgpack:"byte" json:"Byte"`
	ByteSlice       []byte            `msgpack:"byte_slice" json:"ByteSlice"`
	MapStringString map[string]string `msgpack:"map_string_string" json:"MapStringString"`
	MapStringInt    map[string]int    `msgpack:"map_string_int" json:"MapStringInt"`
	MapIntString    map[int]string    `msgpack:"map_int_string" json:"MapIntString"`
}

type TestUnexportedField struct {
	unexportedField int
	ExportedField   int `msgpack:"exported_field"`
}
```

- Ignore a specific field (-i) and ignore a specific struct (-is)

```
st -s -a -v -i=ExportedField -is=TestStructWithTagsCamel -t=msgpack $GOPATH/src/github.com/alistanis/st/etc/etc.go
```

```go
package etc

type TestStruct struct {
	Int             int               `msgpack:"int"`
	Int64           int64             `msgpack:"int_64"`
	IntSlice        []int             `msgpack:"int_slice"`
	Int64Slice      []int64           `msgpack:"int_64_slice"`
	String          string            `msgpack:"string"`
	StringSlice     []string          `msgpack:"string_slice"`
	Float           float64           `msgpack:"float"`
	FloatSlice      []float64         `msgpack:"float_slice"`
	UIntPointer     uintptr           `msgpack:"u_int_pointer"`
	Rune            rune              `msgpack:"rune"`
	RuneSlice       []rune            `msgpack:"rune_slice"`
	Byte            byte              `msgpack:"byte"`
	ByteSlice       []byte            `msgpack:"byte_slice"`
	MapStringString map[string]string `msgpack:"map_string_string"`
	MapStringInt    map[string]int    `msgpack:"map_string_int"`
	MapIntString    map[int]string    `msgpack:"map_int_string"`
}

type TestStructWithTagsSnake struct {
	Int             int               `msgpack:"int" json:"int"`
	Int64           int64             `msgpack:"int_64" json:"int_64"`
	IntSlice        []int             `msgpack:"int_slice" json:"int_slice"`
	Int64Slice      []int64           `msgpack:"int_64_slice" json:"int_64_slice"`
	String          string            `msgpack:"string" json:"string"`
	StringSlice     []string          `msgpack:"string_slice" json:"string_slice"`
	Float           float64           `msgpack:"float" json:"float"`
	FloatSlice      []float64         `msgpack:"float_slice" json:"float_slice"`
	UIntPointer     uintptr           `msgpack:"u_int_pointer" json:"u_int_pointer"`
	Rune            rune              `msgpack:"rune" json:"rune"`
	RuneSlice       []rune            `msgpack:"rune_slice" json:"rune_slice"`
	Byte            byte              `msgpack:"byte" json:"byte"`
	ByteSlice       []byte            `msgpack:"byte_slice" json:"byte_slice"`
	MapStringString map[string]string `msgpack:"map_string_string" json:"map_string_string"`
	MapStringInt    map[string]int    `msgpack:"map_string_int" json:"map_string_int"`
	MapIntString    map[int]string    `msgpack:"map_int_string" json:"map_int_string"`
}

type TestStructWithTagsCamel struct {
	Int             int               `json:"Int"`
	Int64           int64             `json:"Int64"`
	IntSlice        []int             `json:"IntSlice"`
	Int64Slice      []int64           `json:"Int64Slice"`
	String          string            `json:"String"`
	StringSlice     []string          `json:"StringSlice"`
	Float           float64           `json:"Float"`
	FloatSlice      []float64         `json:"FloatSlice"`
	UIntPointer     uintptr           `json:"UIntPointer"`
	Rune            rune              `json:"Rune"`
	RuneSlice       []rune            `json:"RuneSlice"`
	Byte            byte              `json:"Byte"`
	ByteSlice       []byte            `json:"ByteSlice"`
	MapStringString map[string]string `json:"MapStringString"`
	MapStringInt    map[string]int    `json:"MapStringInt"`
	MapIntString    map[int]string    `json:"MapIntString"`
}

type TestUnexportedField struct {
	unexportedField int
	ExportedField   int `msgpack:"-"`
}
```

#Contributing
If you would like to contribute, fork the project, write tests for any new code and ensure that you don't break existing
functionality by running the current tests. If you're looking to submit changes upstream, it would be a good idea to
discuss it with me first. I'm available via email, through Github, or on the Gophers slack team under the st channel 
(https://blog.gopheracademy.com/gophers-slack-community/).

If you do submit a pull request, I will review it and I will merge it if it's in line with my vision for the project.
