package net

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListener(t *testing.T) {
	Convey("Given a servemux", t, func() {
		servemux := ServeMux()
		server := httptest.NewServer(servemux)
		//server.Start()
		str := &StructTagRequest{TagName: "json", Message: string(testStructNoPackageDecl), Case: "snake"}
		requestData, err := json.Marshal(str)
		reader := bytes.NewReader(requestData)
		req, err := http.NewRequest("POST", server.URL+"/tag_struct", reader)
		So(err, ShouldBeNil)
		resp, err := http.DefaultClient.Do(req)
		So(err, ShouldBeNil)
		data, err := ioutil.ReadAll(resp.Body)
		So(err, ShouldBeNil)
		So(string(data), ShouldEqual, expectedStructWithPackageDeclOutput)
		server.Close()
	})
}
