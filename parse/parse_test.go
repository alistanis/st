package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tempDir string

	testDataNoExistingTags = `package test

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
`
	snakeTestDataExistingTags = strings.Replace(`package test

type TestStruct struct {
	Int             int               %sjson:"int"%s
	Int64           int64             %sjson:"int_64"%s
	IntSlice        []int             %sjson:"int_slice"%s
	Int64Slice      []int64           %sjson:"int_64_slice"%s
	String          string            %sjson:"string"%s
	StringSlice     []string          %sjson:"string_slice"%s
	Float           float64           %sjson:"float"%s
	FloatSlice      []float64         %sjson:"float_slice"%s
	UIntPointer     uintptr           %sjson:"u_int_pointer"%s
	Rune            rune              %sjson:"rune"%s
	RuneSlice       []rune            %sjson:"rune_slice"%s
	Byte            byte              %sjson:"byte"%s
	ByteSlice       []byte            %sjson:"byte_slice"%s
	MapStringString map[string]string %sjson:"map_string_string"%s
	MapStringInt    map[string]int    %sjson:"map_string_int"%s
	MapIntString    map[int]string    %sjson:"map_int_string"%s
}
`, "%s", "`", -1)
	camelTestDataExistingTags = strings.Replace(`package test

type TestStruct struct {
	Int             int               %sjson:"Int"%s
	Int64           int64             %sjson:"Int64"%s
	IntSlice        []int             %sjson:"IntSlice"%s
	Int64Slice      []int64           %sjson:"Int64Slice"%s
	String          string            %sjson:"String"%s
	StringSlice     []string          %sjson:"StringSlice"%s
	Float           float64           %sjson:"Float"%s
	FloatSlice      []float64         %sjson:"FloatSlice"%s
	UIntPointer     uintptr           %sjson:"UIntPointer"%s
	Rune            rune              %sjson:"Rune"%s
	RuneSlice       []rune            %sjson:"RuneSlice"%s
	Byte            byte              %sjson:"Byte"%s
	ByteSlice       []byte            %sjson:"ByteSlice"%s
	MapStringString map[string]string %sjson:"MapStringString"%s
	MapStringInt    map[string]int    %sjson:"MapStringInt"%s
	MapIntString    map[int]string    %sjson:"MapIntString"%s
}
`, "%s", "`", -1)

	ignoredStructData = `package test

type TestStruct struct {
	Field string
}
`
	appendData = strings.Replace(`package test

type TestStruct struct {
	Field string %sjson:"field"%s
}
`, "%s", "`", -1)

	expectedAppendData = strings.Replace(`package test

type TestStruct struct {
	Field string %smsgpack:"field" json:"field"%s
}
`, "%s", "`", -1)
)

func init() {
	var err error
	tempDir, err = ioutil.TempDir("", "")
	if err != nil {
		fmt.Println("Could not create temporary directory")
		os.Exit(-1)
	}
}

func Cleanup() {

}

func TestSnakeCase(t *testing.T) {
	Convey("Given sample code with multiple types of structs with tags/no tags", t, func() {
		opts := DefaultOptions()
		SetOptions(opts)
		Convey("We can add the snake case struct tags to it when it has no existing tags", func() {
			data, err := ProcessBytes([]byte(testDataNoExistingTags), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, snakeTestDataExistingTags)
		})
		Convey("We can skip tags if they already exist", func() {
			data, err := ProcessBytes([]byte(camelTestDataExistingTags), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, camelTestDataExistingTags)
		})
		Convey("It will overwrite when overwrite mode is set", func() {
			options.AppendMode = Overwrite
			data, err := ProcessBytes([]byte(camelTestDataExistingTags), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, snakeTestDataExistingTags)
		})

		Convey("It will append a new tag when append mode is set to append", func() {
			opts := DefaultOptions()
			opts.AppendMode = Append
			opts.Tag = "msgpack"
			SetOptions(opts)
			data, err := ProcessBytes([]byte(appendData), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, expectedAppendData)

		})
	})
}

func TestCamelCase(t *testing.T) {
	Convey("Given sample code with multiple types of structs with tags/no tags", t, func() {
		opts := DefaultOptions()
		opts.Case = Camel
		SetOptions(opts)
		Convey("We can add camel case struct tags to it when it has no existing tags", func() {
			data, err := ProcessBytes([]byte(testDataNoExistingTags), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, camelTestDataExistingTags)
		})
		Convey("We can overwrite existing tags when overwrite mode is set", func() {
			options.AppendMode = Overwrite
			data, err := ProcessBytes([]byte(snakeTestDataExistingTags), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, camelTestDataExistingTags)
		})
		Convey("We will leave existing tags alone if they are the same case", func() {
			data, err := ProcessBytes([]byte(camelTestDataExistingTags), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, camelTestDataExistingTags)
		})
	})
}

func TestIgnore(t *testing.T) {
	Convey("Given sample code with multiple types of structs with tags/no tags", t, func() {
		opts := DefaultOptions()
		IgnoredFields = []string{"Field"}
		SetOptions(opts)

		Convey("We can skip an explicitly ignored field", func() {
			data, err := ProcessBytes([]byte(`package test

type TestStruct struct {
	Field string
}
`), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, strings.Replace(`package test

type TestStruct struct {
	Field string %sjson:"-"%s
}
`, "%s", "`", -1))
		})

		Convey("We can skip a struct entirely", func() {
			IgnoredStructs = []string{"TestStruct"}
			data, err := ProcessBytes([]byte(ignoredStructData), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, ignoredStructData)
		})

		Convey("When case isn't set it leaves the struct alone", func() {
			opts := DefaultOptions()
			opts.Case = ""
			opts.Verbose = true
			IgnoredStructs = []string{}
			IgnoredFields = []string{}
			SetOptions(opts)
			data, err := ProcessBytes([]byte(camelTestDataExistingTags), "test.go")
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, camelTestDataExistingTags)
		})

	})
}
func TestFiles(t *testing.T) {
	Convey("Given a temporary directory and temporary files", t, func() {

		Convey("We can process the actual files and the data looks like we expect it to", func() {
			opts := DefaultOptions()
			SetOptions(opts)
			mockFiles := map[string]string{"none": testDataNoExistingTags, "snake": snakeTestDataExistingTags, "camel": camelTestDataExistingTags}
			files := make(map[string]*os.File)

			for k, v := range mockFiles {
				f, err := ioutil.TempFile(tempDir, "")
				So(err, ShouldBeNil)
				files[k] = f
				files[k].WriteString(v)
				err = files[k].Close()
				So(err, ShouldBeNil)
			}
			for k, f := range files {
				data, err := ProcessFile(f.Name())
				So(err, ShouldBeNil)
				switch k {
				case "none":
					So(string(data), ShouldEqual, snakeTestDataExistingTags)
				case "snake":
					So(string(data), ShouldEqual, snakeTestDataExistingTags)
				case "camel":
					So(string(data), ShouldEqual, camelTestDataExistingTags)
				}
			}

			Reset(func() {
				for _, f := range files {
					err := os.Remove(f.Name())
					So(err, ShouldBeNil)
				}
			})
		})

		Convey("Given two files", func() {
			opts := DefaultOptions()
			SetOptions(opts)

			f1, err := ioutil.TempFile(tempDir, "")
			So(err, ShouldBeNil)
			f1.WriteString(testDataNoExistingTags)
			err = f1.Close()
			So(err, ShouldBeNil)
			f2, err := ioutil.TempFile(tempDir, "")
			So(err, ShouldBeNil)
			f2.WriteString(snakeTestDataExistingTags)
			err = f2.Close()
			So(err, ShouldBeNil)

			files := []string{f1.Name(), f2.Name()}
			stdout := os.Stdout
			fname := filepath.Join(tempDir, "stdout")
			temp, err := os.Create(fname)
			So(err, ShouldBeNil)

			os.Stdout = temp

			err = AndProcessFiles(files)
			So(err, ShouldBeNil)

			err = temp.Close()
			So(err, ShouldBeNil)

			output, err := ioutil.ReadFile(fname)
			So(err, ShouldBeNil)

			os.Stdout = stdout

			So(strings.Trim(string(output), "\n")+"\n", ShouldContainSubstring, strings.Trim(snakeTestDataExistingTags+"\n"+snakeTestDataExistingTags, "\n"))

			Reset(func() {
				for _, f := range files {
					err = os.Remove(f)
					So(err, ShouldBeNil)
				}
			})

		})

		Convey("Given a set of 'files' (as defined in parse.go)", func() {
			files := []*File{{FileName: "test.go", Data: []byte(testDataNoExistingTags)}}
			results, err := Process(files)
			So(err, ShouldBeNil)
			So(string(results[0].Data), ShouldEqual, snakeTestDataExistingTags)
		})

		Convey("Given a src file", func() {
			f, err := ioutil.TempFile(tempDir, "")
			So(err, ShouldBeNil)
			_, err = f.WriteString(testDataNoExistingTags)
			So(err, ShouldBeNil)
			opts := DefaultOptions()
			opts.DryRun = false
			SetOptions(opts)
			err = AndProcessFiles([]string{f.Name()})
			So(err, ShouldBeNil)
			data, err := ioutil.ReadFile(f.Name())
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, snakeTestDataExistingTags)

			Reset(func() {
				err = os.Remove(f.Name())
				So(err, ShouldBeNil)
			})
		})

	})
}

func TestErrors(t *testing.T) {
	Convey("Given a malformed string of code", t, func() {
		badSrc := `package test
asofinqowkernoaskn{{}}'`
		Convey("An error is returned immediately", func() {
			data, err := ProcessBytes([]byte(badSrc), "test.go")
			So(data, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given a bad file path", t, func() {
		badPath := ""
		Convey("An error is returned immediately", func() {
			data, err := ProcessFile(badPath)
			So(data, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given bad files and directories", t, func() {
		Convey("Empty string in a path of files returns an error", func() {
			err := AndProcessFiles([]string{""})
			So(err, ShouldNotBeNil)
		})

		Convey("Passing a directory returns an error", func() {
			err := AndProcessFiles([]string{tempDir})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, fmt.Sprintf("Cannot use a directory as a path."))
		})
	})

	Convey("Given a struct with no fields, we can continue silently and will not return an error", t, func() {
		emptyStruct := `package test

type Embedded struct{}

type TestStruct struct {
	Embedded
}
`
		data, err := ProcessBytes([]byte(emptyStruct), "test.go")
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, emptyStruct)
	})

	Convey("Given a list of files that have bad source, we can return the last error that occurs", t, func() {
		badSrc1 := `package test
		aposmpomparosmpckmqlkwme`
		badSrc2 := `package test
		asdpomqpwoermpoamspdlmdw`
		files := []*File{{FileName: "bad_1", Data: []byte(badSrc1)}, {FileName: "bad_2", Data: []byte(badSrc2)}}
		_, err := Process(files)
		So(err.Error(), ShouldContainSubstring, "bad_2")
	})
}

// Below taken from https://github.com/etgryphon/stringUp/blob/master/stringUp_test.go
func TestCamelCased(t *testing.T) {
	const cameled, upCameled = "thisIsIt", "ThisIsItBob"
	if x := CamelCase(cameled); x != cameled {
		t.Errorf("CamelCase(%v) = %v, want %v", cameled, x, cameled)
	}
	if x := CamelCase(upCameled); x != upCameled {
		t.Errorf("CamelCase(%v) = %v, want %v", upCameled, x, upCameled)
	}
}

func TestCamelCaseSpaced(t *testing.T) {
	const src, upSrc = "this is it", "This Is It Bob"
	if x := CamelCase(src); x != "thisIsIt" {
		t.Errorf("CamelCase(%v) = %v, want %v", src, x, "thisIsIt")
	}
	if x := CamelCase(upSrc); x != "ThisIsItBob" {
		t.Errorf("CamelCase(%v) = %v, want %v", upSrc, x, "ThisIsItBob")
	}
}

func TestCamelCaseUnderscored(t *testing.T) {
	const src, upSrc = "this_is_it", "This_Is_It_Bob"
	if x := CamelCase(src); x != "thisIsIt" {
		t.Errorf("CamelCase(%v) = %v, want %v", src, x, "thisIsIt")
	}
	if x := CamelCase(upSrc); x != "ThisIsItBob" {
		t.Errorf("CamelCase(%v) = %v, want %v", upSrc, x, "ThisIsItBob")
	}
}

func TestCamelCaseDashed(t *testing.T) {
	const src, upSrc = "this-is-it", "This-Is-It-Bob"
	if x := CamelCase(src); x != "thisIsIt" {
		t.Errorf("CamelCase(%v) = %v, want %v", src, x, "thisIsIt")
	}
	if x := CamelCase(upSrc); x != "ThisIsItBob" {
		t.Errorf("CamelCase(%v) = %v, want %v", upSrc, x, "ThisIsItBob")
	}
}

func TestCamelCaseMixed(t *testing.T) {
	const src, upSrc = "-this is_it", "This Is_It-Bob"
	if x := CamelCase(src); x != "thisIsIt" {
		t.Errorf("CamelCase(%v) = %v, want %v", src, x, "thisIsIt")
	}
	if x := CamelCase(upSrc); x != "ThisIsItBob" {
		t.Errorf("CamelCase(%v) = %v, want %v", upSrc, x, "ThisIsItBob")
	}
}
