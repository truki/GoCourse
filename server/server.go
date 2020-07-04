package main

import (
	"github.com/pabloos/http/greet"
	"log"
	"net/http"
	"os"
)

func newServer(cache map[string]greet.Greet) *http.Server {
	return &http.Server{
		Addr:      ":8080",
		Handler:   newMux(cache),
		TLSConfig: tlsConfig(),
		ErrorLog:  log.New(os.Stderr, "HTTP Server says: ", log.Llongfile),
	}
}
