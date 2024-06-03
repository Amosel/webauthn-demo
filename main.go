package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var webAuthnSessionStore *WebAuthnSessionStore

func main() {
	log.Println("Initializing the WebAuthn session store")
	webAuthnSessionStore = NewWebAuthnSessionStore()

	log.Println("Setting up router and handlers")
	r := mux.NewRouter()
	r.HandleFunc("/webauthn/register/begin", beginWebAuthnRegistrationHandler).Methods("POST")
	r.HandleFunc("/webauthn/register/finish", finishWebAuthnRegistrationHandler).Methods("POST")
	r.HandleFunc("/webauthn/login/begin", beginWebAuthnLoginHandler).Methods("POST")
	r.HandleFunc("/webauthn/login/finish", finishWebAuthnLoginHandler).Methods("POST")

	// Serve static files from the "static" directory
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Println("Starting HTTPS server on port 443")
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
