package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alistanis/st/sterrors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestListener(t *testing.T) {
	Convey("Given a servemux", t, func() {
		servemux := ServeMux()
		server := httptest.NewServer(servemux)
		//server.Start()
		str := &StructTagRequest{TagName: "json", Message: string(testStructNoPackageDecl), Case: "snake"}
		requestData, err := json.Marshal(str)
		So(err, ShouldBeNil)
		reader := bytes.NewReader(requestData)
		Convey("we can make a request to the test server and receive the expected response", func() {
			req, err := http.NewRequest("POST", server.URL+"/tag_struct", reader)
			So(err, ShouldBeNil)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			data, err := ioutil.ReadAll(resp.Body)
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, expectedStructWithPackageDeclOutput)
		})
		Convey("We can make a request to the test server with invalid syntax and receive the expected error", func() {
			req, err := http.NewRequest("POST", server.URL+"/tag_struct", bytes.NewReader([]byte("BAD SYNTAX")))
			So(err, ShouldBeNil)
			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			data, err := ioutil.ReadAll(resp.Body)
			So(err, ShouldBeNil)
			fmt.Println(string(data))
			httpErr := &sterrors.HttpError{}
			err = json.Unmarshal(data, &httpErr)
			So(err, ShouldBeNil)
			So(httpErr.Code, ShouldEqual, 400)
		})

		server.Close()
	})

	// This is always a little trickier to test, because ServeHTTP is a blocking function we spin it off into another goroutine
	// Then we create a boolean channel to block at the end of the Convey scope which we will send to at the end of second goroutine
	// that possesses its own GoConvey context. This ensures that the server will run until the second goroutine resturns
	Convey("Given an actual http server running in a goroutine", t, func() {
		go func() {
			ServeHTTP(ServeMux())
		}()

		block := make(chan bool)

		go func() {
			Convey("We can make a request from a separate goroutine", t, func(c C) {
				str := &StructTagRequest{TagName: "json", Message: string(testStructNoPackageDecl), Case: "snake"}
				requestData, err := json.Marshal(str)
				req, err := http.NewRequest("POST", "http://localhost:8080/tag_struct", bytes.NewReader(requestData))
				c.So(err, ShouldBeNil)
				resp, err := http.DefaultClient.Do(req)
				c.So(err, ShouldBeNil)
				data, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(string(data), ShouldEqual, expectedStructWithPackageDeclOutput)
			})
			block <- true
		}()

		<-block
	})
}
