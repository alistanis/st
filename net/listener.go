package net

import (
	"net/http"

	"github.com/alistanis/st/sterrors"
)

func Serve() {
	servemux := http.NewServeMux()
	servemux.HandleFunc("/tag_struct", func(rw http.ResponseWriter, req *http.Request) {
		resp, err := processStructTagRequest(req)
		if err != nil {
			rw.WriteHeader(400)
			rw.Write(sterrors.FormatHTTPError(err, 400))
			return
		}
		rw.WriteHeader(200)
		rw.Write(resp)
	})
	http.ListenAndServe(":8080", servemux)
}
