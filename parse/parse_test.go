package parse

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testDataNoExistingTags = `type TestStruct struct {
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
}`
	snakeTestDataExistingTags = strings.Replace(`type TestStructWithTagsSnake struct {
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
}`, "%s", "`", -1)
	camelTestDataExistingTags = strings.Replace(`type TestStructWithTagsCamel struct {
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
}`, "%s", "`", -1)
)

func TestSnakeCase(t *testing.T) {
	Convey("Given a sample piece of code defining a struct with no struct tags and the correct options", t, func() {
		opts := DefaultOptions
		SetOptions(opts)
		Convey("We can add the appropriate struct tags to it", func() {
			data, err := ProcessBytes([]byte(testDataNoExistingTags), "test.go")
			So(err, ShouldBeNil)
			So(data, ShouldEqual, []byte(snakeTestDataExistingTags))
		})
	})
}
