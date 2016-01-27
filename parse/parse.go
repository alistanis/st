// Package parse provides functions for parsing and tagging golang structs. It achieves this by creating an ast and
// visiting all of its nodes.
package parse

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"reflect"

	"go/format"

	"path/filepath"

	"github.com/alistanis/st/sterrors"
)

// Append modes
const (
	// Append will append to tags rather than overwriting them altogether
	Append = iota
	// Overwrite will overwrite tags completely
	Overwrite
	// SkipExisting skips existing tags whether or not they match tag or case
	SkipExisting
)

// Major Tag modes
const (
	// TagAll will tag all structs/fields (unless they are excluded in the IngoreStructs/IgnoreFields slices)
	TagAll = iota
	// SkipSpecifiedStructs (currently unimplemented)
	SkipSpecifiedStructs
	// IncludeSpecifiedStructs (currently unimplemented)
	IncludeSpecifiedStructs
	// SkipStructAndFieldKeypairs (currently unimplemented)
	SkipStructAndFieldKeypairs
	// IncludeStructAndFieldKeypairs (currently unimplemented)
	IncludeStructAndFieldKeypairs
)

// Basic supported tags and cases
const (
	// JSON represents the json tag
	JSON = "json"
	// Snake represents snake case
	Snake = "snake"
	// Camel represents camel case
	Camel = "camel"
	// DefaultGenerateTag represents the default go generate tag that ST will respect
	DefaultGenerateTag = "@st"
)

// Defaults
var (
	// DefaultAppendMode is SkipExisting - will skip existing tags entirely
	DefaultAppendMode = SkipExisting
	// DefaultTagMode is TagAll - will tag all structs/fields unless they are already tagged or in the excluded slices
	DefaultTagMode = TagAll
	// DefaultTag is JSON
	DefaultTag = JSON
	// DefaultCase is Snake case. (common in http, sql, etc)
	DefaultCase = Snake

	options = DefaultOptions()
	// IgnoredFields contains strings for fields that are not to be tagged
	IgnoredFields = make([]string, 0)
	// IgnoredStructs contains strings for structs that are not to be tagged
	IgnoredStructs = make([]string, 0)
	// TODO - Add more sophisticated exclusion/inclusion after refactor

)

// localGlobals
var (
	lastCommentWithGenerateTag string
	lastTypeName               string
)

// CommentDirective represents a comment with //@st at its beginning.
// I am really not a fan of treating comments as anything more than a comment, but Go unfortunately has no other constructs
type CommentDirective struct {
	BaseText string
	FlagSet  *flag.FlagSet
}

// Args returns the underlying FlagSet.Args()
func (c *CommentDirective) Args() []string {
	return c.FlagSet.Args()
}

/* NewCommentDirective takes a string (which should be the text from an *ast.Comment), creates a new flag set using the comment
   as the flag set name with flag.ContinueOnError - flag.ExitOnError will call os.Exit() - and then parses the flags
   Comments: 1) I would pass in the *ast.Comment directly, but it is already initialized further up the call stack at this point
   			 2) There is some hackery of the flag package going on in here
   s should be in the following format:
   - It should only be one line
   - It should say //@st with no spaces
   - Commands given after //@st will be interpreted just like normal st commands
   - $GOFILE will be passed in as the final argument */
func NewCommentDirective(s string) (*CommentDirective, error) {
	cd := &CommentDirective{BaseText: s, FlagSet: flag.NewFlagSet(s, flag.ContinueOnError)}
	args := strings.Split(strings.TrimLeft(s, `//@st`), " ")
	SetArgs(append(args, os.Getenv("GOFILE")))
	err := cd.FlagSet.Parse(args)
	return cd, err
}

// DefaultOptions returns a new *Options with all default values initialized
func DefaultOptions() *Options {
	return &Options{
		//Tags:       []string{DefaultTag},
		Tag:         DefaultTag,
		Case:        DefaultCase,
		AppendMode:  DefaultAppendMode,
		TagMode:     DefaultTagMode,
		DryRun:      true,
		Verbose:     false,
		GenerateTag: DefaultGenerateTag}
}

// SetOptions sets the current options to the options provided. (This is not thread safe if called from a goroutine)
func SetOptions(o *Options) {
	options = o
}

// Options represents package behavior options - will be expanded to take a list of tags to support go generate
type Options struct {
	//Tags       []string
	Tag         string
	Case        string
	AppendMode  int
	TagMode     int
	DryRun      bool
	Verbose     bool
	GenerateTag string
}

// AndProcessFiles takes a list of paths, iterates over them, stats them, and then inspects source files
func AndProcessFiles(paths []string) error {
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return fmt.Errorf("Cannot use a directory as a path. Path: %s", fi.Name())
		}
		data, err := ProcessFile(p)
		if options.DryRun {
			fmt.Println(string(data))
		} else {
			ioutil.WriteFile(p, data, 0664)
		}
	}
	return nil
}

// File represents a basic file with a FileName(path) and the Data contained within the file
type File struct {
	FileName string
	Data     []byte
}

// Process iterates over a []*File, processes the *Files, and returns the resulting []*File and the last error that occurred, if any
// This function could potentially consume a lot of memory if an extraordinarily large set was passed to it
func Process(files []*File) ([]*File, error) {
	var lastErr error
	var results []*File
	for _, f := range files {
		data, err := ProcessBytes(f.Data, f.FileName)
		if err != nil {
			lastErr = err
			continue
		}
		result := &File{FileName: f.FileName, Data: data}
		results = append(results, result)
	}
	return results, lastErr
}

// ProcessBytes takes a []byte and filename, and inspects the data, returning that data in another []byte
func ProcessBytes(data []byte, filename string) ([]byte, error) {
	astFile, data, err := Parse(data, filename)
	if err != nil {
		return nil, err
	}
	return Inspect(astFile, data)
}

// Parse returns an *ast.File, the data parsed, and an error
func Parse(data []byte, filename string) (*ast.File, []byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, string(data), parser.ParseComments)
	return f, data, err
}

// ProcessFile processes a file, returning the processed []byte
func ProcessFile(path string) ([]byte, error) {
	f, data, err := parseFile(path)
	if err != nil {
		return nil, err
	}
	return Inspect(f, data)
}

// parseFile reads all file information into a buffer, then creates a token set and parses the file, returning a *ast.File
func parseFile(path string) (*ast.File, []byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	name := filepath.Base(path)
	return Parse(data, name)
}

// Inspect visits all nodes in the *ast.File (recursively), performing mutations on the buffer when the type found is an *ast.StructType
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
		case *ast.Comment:
			if strings.Contains(t.Text, options.GenerateTag) {
				lastCommentWithGenerateTag = strings.TrimLeft(t.Text, `//`)
			}
		case *ast.StructType:
			data = TagStruct(data, t, offset)
		}
		return true
	})
	return format.Source(data)
}

// TagStruct tags a struct based on whether or not it is exported, is ignored, and what flags are provided at runtime
func TagStruct(srcData []byte, s *ast.StructType, offset *int) []byte {
	// If the last type name is one of our ignored structs, return immediately
	if IsIgnoredTypeName(lastTypeName) {
		return srcData
	}
	for _, f := range s.Fields.List {
		if len(f.Names) == 0 {
			sterrors.Printf("Could not find name for field: %+v\n", f)
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
				if options.AppendMode == SkipExisting || options.AppendMode == Append {
					currentTagValue := reflectTag.Get(options.Tag)
					if currentTagValue != "" {
						sterrors.Printf("Existing tag found: TagName: %s, TagValue: %s, StartIndex: %d, EndIndex: %d - Skipping Tag\n", options.Tag, currentTagValue, tag.Pos(), tag.End())
						continue
					}
				}
				srcData = OverwriteStructTag(tag, formattedName, offset, srcData)
			} else {
				srcData = AppendStructTag(f, formattedName, offset, srcData)
			}
		}

	}
	return srcData
}

// AppendStructTag adds an additional tag to a struct tag
func AppendStructTag(field *ast.Field, tagName string, offset *int, data []byte) []byte {
	start := int(field.End()) + *offset - 1
	tag := fmt.Sprintf(" `%s:\"%s\"`", options.Tag, tagName)
	*offset += len(tag)
	return Insert(data, []byte(tag), start)
}

// OverwriteStructTag overwrites the struct tag completely
func OverwriteStructTag(tag *ast.BasicLit, tagName string, offset *int, data []byte) []byte {
	val := tag.Value
	start := int(tag.Pos()) + *offset - 1
	end := int(tag.End()) + *offset - 1
	length := len(val)
	oldLength := end - start

	// Delete the original tag
	data = DeleteRange(data, start, end)
	var newTag string
	if options.AppendMode == Append {
		oldTag := removeIndex(removeIndex(val, 0), len(val)-2)
		newTag = fmt.Sprintf("`%s:\"%s\" %s`", options.Tag, tagName, oldTag)
	} else {
		newTag = fmt.Sprintf("`%s:\"%s\"`", options.Tag, tagName)
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

// IsIgnoredField checks if a field is an explicitly ignored field
// Currently a slice is fine for performance, but we will replace these with maps later.
func IsIgnoredField(s string) bool {
	for _, f := range IgnoredFields {
		if s == f {
			return true
		}
	}
	return false
}

// IsIgnoredTypeName checks if the name provided is an ignored struct
func IsIgnoredTypeName(s string) bool {
	for _, n := range IgnoredStructs {
		if n == s {
			return true
		}
	}
	return false
}

// DeleteRange deletes a range from a []byte, returning the new slice
func DeleteRange(data []byte, start, end int) []byte {
	return append(data[:start], data[end:]...)
}

// Insert inserts insertData into data at the given start index
func Insert(data, insertData []byte, start int) []byte {
	return append(data[:start], append(insertData, data[start:]...)...)
}

// removeIndex removes a single index from a string
func removeIndex(input string, index int) string {
	return input[:index] + input[index+1:]
}

// FormatFieldName formats the field name as either CamelCase or snake_case
func FormatFieldName(n string) string {
	switch options.Case {
	case Camel:
		return CamelCase(n)
	case Snake:
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

// Underscore will change a string from a camelcased
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

// Taken from https://github.com/etgryphon/stringUp/blob/master/stringUp.go
var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

// CamelCase converts a string to the CamelCase version of it
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
