package net

import "net/http"

func ServeStaticContent() error {
	servemux := http.NewServeMux()
	servemux.Handle("/", http.FileServer(http.Dir("./public")))
	return http.ListenAndServe("localhost:8085", servemux)
}
