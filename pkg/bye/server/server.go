package server

import (
	"fmt"
	"net/http"
)

// RunServer starts the server on the given port.
func Run(dir string, port int64) {
	if port == 0 {
		port = 3005
	}

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", http.StripPrefix("/", fs))

	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}
