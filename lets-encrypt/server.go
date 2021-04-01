package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

// User implements acme.User
type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}
func (u User) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func main() {
	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	brian := User{
		Email: "brian@email.com",
		key:   privateKey,
	}

	config := lego.NewConfig(&brian)

	config.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// start up the server for the HTTP01 Challenge
	providerServer := http01.NewProviderServer("", "")
	err = client.Challenge.SetHTTP01Provider(providerServer)
	if err != nil {
		log.Fatal(err)
	}

	// New users will need to register
	opts := registration.RegisterOptions{TermsOfServiceAgreed: true}
	reg, err := client.Registration.Register(opts)
	if err != nil {
		log.Fatal(err)
	}
	brian.Registration = reg

	// request certs for 1.foshee.dev domain
	request := certificate.ObtainRequest{
		Domains: []string{"2.foshee.dev"},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	// cert is bundled with issuer cert, no need to store both
	cert := certificates.Certificate

	// save the cert
	f, err := os.Create("1.crt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(f, "%s\n", cert)
	f.Close()

	pk := certificates.PrivateKey

	// save the private key
	pkf, err := os.Create("1.key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(pkf, "%s\n", pk)
	pkf.Close()

	// challenge server is shut down at this point
	// start up a server to handle app requests
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from 1.foshee.dev")
	})

	server := http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}

	fmt.Printf("Server listening on %s", server.Addr)
	if err := server.ListenAndServeTLS("1.crt", "1.key"); err != nil {
		fmt.Println(err)
	}
}
