package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pabloos/http/greet"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
	"strings"
)

func POST(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func Debug(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer h.ServeHTTP(w, r)

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(dump))
	}
}

func Delay(delay time.Duration, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer h.ServeHTTP(w, r)

		time.Sleep(delay)
	}
}

func caching(cache map[string]greet.Greet, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		// rewriting body of request
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		defer h.ServeHTTP(w, r)
		var t greet.Greet
		err = json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		// rewitting body of request
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		val, ok := cache[t.Name]
		if ok {
			log.Printf("Resource %v in cache:", val)
		} else {
			log.Printf("Resource %v is not cached", t.Name)
			cache[t.Name] = t

			var similars = make([]greet.Greet,0)
			for _, val := range cache {
				if strings.Contains(val.Name, t.Name) && (val.Name != t.Name) {
					similars = append(similars, val)
				}
			}
			log.Println("Recursos similares: ", similars)
			if len(similars) > 0 {
				fmt.Fprintf(w, "Los siguientes recursos son similares: %v\n", similars)
			}


		}
		log.Println("CACHE Now: ", cache)
	}
}