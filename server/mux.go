package main

import (
	"github.com/pabloos/http/greet"
	"net/http"
)

func newMux(cache map[string]greet.Greet) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Debug(index))
	//mux.HandleFunc("/greet", Delay(2*time.Second, POST(greetHandler)))
	mux.HandleFunc("/greet", caching(cache, greetHandler))

	return mux
}
