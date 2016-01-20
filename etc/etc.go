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

type TestUnexported struct {
	testUnexportedInt int
	TestExportedInt   int
}
