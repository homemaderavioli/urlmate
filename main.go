package main

import (
	"fmt"
	"io"
	"log"
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	srv := &server{
		domain: "localhost",
		port:   port,
		db: &RiakClient{
			IP:             "localhost",
			ShortURLBucket: "short_urls",
		},
		router: http.NewServeMux(),
	}
	srv.routes()

	log.Printf("Listening on port %s", srv.port)
	return http.ListenAndServe(":"+srv.port, srv)
}
