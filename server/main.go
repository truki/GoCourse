package main

import (
	"github.com/pabloos/http/greet"
	"log"
)

func main() {
	var cache = make(map[string]greet.Greet)
	if err := newServer(cache).ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
