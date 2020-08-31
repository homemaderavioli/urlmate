package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	srv := &server{
		domain: "localhost",
		port:   ":8080",
		db: &RiakClient{
			IP:             "localhost",
			ShortURLBucket: "short_urls",
		},
		router: http.NewServeMux(),
	}
	srv.routes()
	return http.ListenAndServe(srv.port, srv)
}
