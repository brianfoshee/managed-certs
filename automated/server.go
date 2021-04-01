package main

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello MAPS!")
	})

	m := &autocert.Manager{
		Cache:      autocert.DirCache("certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("maps.foshee.dev"),
	}
	s := &http.Server{
		Addr:      ":https",
		Handler:   mux,
		TLSConfig: m.TLSConfig(),
	}
	s.ListenAndServeTLS("", "")
}
