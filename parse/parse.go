// Most functions in package parse are exported on the off change that someone would like to use them as library functions
// in their own project
package parse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"reflect"

	"flag"

	"go/format"

	"github.com/alistanis/st/flags"
	"github.com/alistanis/st/sterrors"
)

var (
	lastTypeName string
)

// Takes all provided arguments, iterates over them, stats them, and then inspects source files
func ParseAndProcess() error {
	for _, p := range flag.Args() {
		fi, err := os.Stat(p)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("Cannot use a directory as a path. Path: %s", fi.Name())
		}
		f, data, err := ParseFile(p)
		if err != nil {
			return err
		}
		data, err = Inspect(f, data)
		if err != nil {
			return err
		}
		if flags.Write {
			ioutil.WriteFile(p, data, 0664)
		} else {
			fmt.Println(string(data))
		}
	}
	return nil
}

// Reads all file information into a buffer, then creates a token set and parses the file, returning a *ast.File
func ParseFile(path string) (*ast.File, []byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "etc.go", string(data), 0)
	return f, data, err
}

// Visits all nodes in the *ast.File (recursively), performing mutations on the buffer when the type found is an *ast.StructType
func Inspect(f *ast.File, srcFileData []byte) ([]byte, error) {
	data := srcFileData
	var offset *int
	offsetVal := 0
	offset = &offsetVal
	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.Ident:
			if t.Obj != nil {
				if t.Obj.Kind == ast.Typ {
					lastTypeName = t.Obj.Name
				}
			}
		case *ast.StructType:
			data = TagStruct(data, t, offset)
		}
		return true
	})
	return format.Source(data)
}

// Tags a struct based on whether or not it is exported, is ignored, and what flags are provided at runtime
func TagStruct(srcData []byte, s *ast.StructType, offset *int) []byte {
	// If the last type name is one of our ignored structs, return immediately
	if IsIgnoredTypeName(lastTypeName) {
		return srcData
	}
	for _, f := range s.Fields.List {
		if len(f.Names) == 0 {
			fmt.Printf("Could not find name for field: %+v\n", f)
			continue
		}
		if f.Names[0].IsExported() {
			name := f.Names[0].Name
			var formattedName string
			if IsIgnoredField(name) {
				formattedName = "-"
			} else {
				formattedName = FormatFieldName(name)
			}
			tag := f.Tag
			if tag != nil {
				val := tag.Value
				// remove `'s from string and convert to a reflect.StructTag so we can use reflect.StructTag().Get() call
				reflectTag := reflect.StructTag(val[1 : len(val)-1])
				if !flags.Overwrite {
					currentTagValue := reflectTag.Get(flags.Tag)
					if currentTagValue != "" {
						sterrors.Printf("Existing tag found: TagName: %s, TagValue: %s, StartIndex: %d, EndIndex: %d - Skipping Tag\n", flags.Tag, currentTagValue, tag.Pos(), tag.End())
						continue
					}
				}
				srcData = OverwriteStructTag(tag, formattedName, offset, srcData)
			} else {
				srcData = AddStructTag(f, formattedName, offset, srcData)
			}
		}

	}
	return srcData
}

// Adds an additional tag to a struct tag
func AddStructTag(field *ast.Field, tagName string, offset *int, data []byte) []byte {
	start := int(field.End()) + *offset - 1
	tag := fmt.Sprintf(" `%s:\"%s\"`", flags.Tag, tagName)
	*offset += len(tag)
	return Insert(data, []byte(tag), start)
}

// Overwrites the struct tag completely
func OverwriteStructTag(tag *ast.BasicLit, tagName string, offset *int, data []byte) []byte {
	val := tag.Value
	start := int(tag.Pos()) + *offset - 1
	end := int(tag.End()) + *offset - 1
	length := len(val)
	oldLength := end - start

	// Delete the original tag
	data = DeleteRange(data, start, end)
	var newTag string
	if flags.Append {
		oldTag := removeIndex(removeIndex(val, 0), len(val)-2)
		newTag = fmt.Sprintf("`%s:\"%s\" %s`", flags.Tag, tagName, oldTag)
	} else {
		newTag = fmt.Sprintf("`%s:\"%s\"`", flags.Tag, tagName)
	}

	numSpaces := len(newTag) - oldLength - 1
	var spaces string

	// Can't pass a negative number to strings.Repeat()
	// it will cause a panic because it passes this number directly to make()
	if numSpaces > 0 {
		spaces = strings.Repeat(" ", numSpaces)
	}

	newTag = fmt.Sprintf("%s%s", spaces, newTag)
	localOffset := len(newTag) - length
	*offset += localOffset

	// Insert new tag
	data = Insert(data, []byte(newTag), start)
	return data
}

// Checks if a fiels is an explicitly ignored field
func IsIgnoredField(s string) bool {
	for _, f := range flags.IgnoredFields {
		if s == f {
			return true
		}
	}
	return false
}

func IsIgnoredTypeName(s string) bool {
	for _, n := range flags.IgnoredStructs {
		if n == s {
			return true
		}
	}
	return false
}

// Deletes a range from a slice, returning the new slice
func DeleteRange(data []byte, start, end int) []byte {
	return append(data[:start], data[end:]...)
}

// Inserts []byte at the given start index
func Insert(data, insertData []byte, start int) []byte {
	return append(data[:start], append(insertData, data[start:]...)...)
}

// Removes a single index from a string
func removeIndex(input string, index int) string {
	return input[:index] + input[index+1:]
}

// Formats the field name as either CamelCase or snake_case
func FormatFieldName(n string) string {
	switch flags.Case {
	case flags.Camel:
		return CamelCase(n)
	case flags.Snake:
		return Underscore(n)
	}
	sterrors.Printf("Could not format string, Case is not set.\n")
	return n
}

var (
	doubleColon       = regexp.MustCompile("::")
	dash              = regexp.MustCompile("-")
	uppersOrNumsLower = regexp.MustCompile("([A-Z0-9]+)([A-Z][a-z])")
	lowerUpper        = regexp.MustCompile("([a-z])([A-Z0-9])")
)

// This function will change a string from a camelcased
// form to a string with underscores. Will change "::" to
// "/" to maintain compatibility with Rails's underscore
func Underscore(str string) string {
	output := doubleColon.ReplaceAllString(str, "/")

	// Rails uses underscores while I use dashes in this function
	// Go's regexp doesn't like $1_$2, so we'll use a dash instead
	// since it will get fixed in a later replacement
	output = uppersOrNumsLower.ReplaceAllString(output, "$1-$2")
	output = lowerUpper.ReplaceAllString(output, "$1-$2")

	output = strings.ToLower(output)
	output = dash.ReplaceAllString(output, "_")

	return output
}

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

// Converts a string to the CamelCase version of it
func CamelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}
