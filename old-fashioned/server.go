package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello *secure localhost")
	})

	server := http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}

	fmt.Printf("Server listening on %s", server.Addr)

	// generate these with:
	// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout localhost.key -out localhost.crt
	if err := server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		fmt.Println(err)
	}
}
