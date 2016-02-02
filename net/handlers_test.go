package net

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testStructWithPackageDecl = `package st
type Struct struct {
	Field string
}
`

	testStructNoPackageDecl = `type Struct struct {
		Field string
}
`
)

var (
	expectedStructWithPackageDeclOutput = strings.Replace(`package st

type Struct struct {
	Field string %sjson:"field"%s
}
`, "%s", "`", -1)
	expectedStructNoPackageDecl = strings.Replace(`package st

type Struct struct {
	Field string %sjson:"field"%s
}
`, "%s", "`", -1)
)

func TestStructTagHandler(t *testing.T) {
	Convey("Given a set of fake requests", t, func() {

		Convey("We can send a request with a valid package declaration to the handler and receive an appropriate response", func() {
			str := &StructTagRequest{TagName: "json", Message: string(testStructWithPackageDecl), Case: "snake"}
			requestData, err := json.Marshal(str)
			So(err, ShouldBeNil)
			buffer := bytes.NewReader(requestData)
			req, err := http.NewRequest("POST", "http://localhost:8080", buffer)
			So(err, ShouldBeNil)
			data, err := processStructTagRequest(req)
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, expectedStructWithPackageDeclOutput)
		})
		Convey("We can send a request with no package declaration to the handler and receive an appropriate response", func() {
			str := &StructTagRequest{TagName: "json", Message: string(testStructNoPackageDecl), Case: "snake"}
			requestData, err := json.Marshal(str)
			So(err, ShouldBeNil)
			buffer := bytes.NewReader(requestData)
			req, err := http.NewRequest("POST", "http://localhost:8080", buffer)
			So(err, ShouldBeNil)
			data, err := processStructTagRequest(req)
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, expectedStructNoPackageDecl)
		})

		Convey("We can send an invalid request with no body", func() {
			req, err := http.NewRequest("POST", "http://localhost:8080", nil)
			So(err, ShouldBeNil)
			data, err := processStructTagRequest(req)
			So(err, ShouldNotBeNil)
			So(data, ShouldBeNil)
		})

		Convey("We can send an invalid request with plain text in the body", func() {
			req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewReader([]byte("Plain text, no json here")))
			So(err, ShouldBeNil)
			data, err := processStructTagRequest(req)
			So(err, ShouldNotBeNil)
			So(data, ShouldBeNil)
		})
	})
}
