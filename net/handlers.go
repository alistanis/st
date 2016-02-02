package net

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/alistanis/st/parse"
)

type StructTagRequest struct {
	Message    string `json:"message"`
	AppendMode string `json:"append_mode"`
	TagMode    string `json:"tag_mode"`
	TagName    string `json:"tag_name"`
	Case       string `json:"case"`
}

func processStructTagRequest(req *http.Request) (data []byte, err error) {
	// in order to catch a nil request body
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		var ok bool
		if err, ok = e.(error); ok {
			if err.Error() == "runtime error: invalid memory address or nil pointer dereference" {
				err = errors.New("Empty request body. Must send request as json.")
				return
			}
		}
	}()

	var body []byte
	body, _ = ioutil.ReadAll(req.Body)
	str := &StructTagRequest{}
	err = json.Unmarshal(body, &str)
	if err != nil {
		return nil, err
	}
	opts := parse.DefaultOptions()
	opts.Case = str.Case
	opts.Tag = str.TagName
	//not thread safe, will need a context in the future
	parse.SetOptions(opts)
	data = []byte(str.Message)
	// naive implementation for first version, we should use the ast in the final version.
	if !strings.Contains(string(data), "package") {
		data = parse.Insert(data, []byte("package st\n"), 0)
	}
	return parse.ProcessBytes(data, "st.go")
}
