package net

import (
	"net/http"

	"github.com/alistanis/st/sterrors"
)

func ServeMux() *http.ServeMux {
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
	return servemux
}

func ServeHTTP(mux *http.ServeMux) {
	http.ListenAndServe(":8080", mux)
}
