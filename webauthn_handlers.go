package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

func AddCredentialToUser(ctx context.Context, email string, credential *webauthn.Credential) error {
	_, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	// Here you would save the credential to your database
	// This is a placeholder for the actual implementation.
	return nil
}

func beginWebAuthnRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || !isValidEmail(req.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	json.NewDecoder(r.Body).Decode(&req)

	user := User{Email: req.Email, DisplayName: "User Display Name"}
	config := WebAuthnConfiguration{RPID: "example.com", RPOrigin: "https://example.com"}

	log.Println("Starting WebAuthn registration")
	options, err := webAuthnSessionStore.BeginWebAuthnRegistration(user, nil, config)
	if err != nil {
		log.Printf("Error during WebAuthn registration: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("WebAuthn registration options sent")
	json.NewEncoder(w).Encode(options)
}

func finishWebAuthnRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string
		Response string
	}
	json.NewDecoder(r.Body).Decode(&req)

	user := User{Email: req.Email}
	config := WebAuthnConfiguration{RPID: "example.com", RPOrigin: "https://example.com"}

	credential, err := webAuthnSessionStore.FinishWebAuthnRegistration(user, nil, r, config)
	if err != nil {
		log.Printf("Error finishing WebAuthn registration: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("WebAuthn registration completed successfully")

	err = AddCredentialToUser(context.Background(), req.Email, credential)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func beginWebAuthnLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string
	}
	json.NewDecoder(r.Body).Decode(&req)

	user := User{Email: req.Email}
	sr := SessionRequest{
		WebAuthnConfig: WebAuthnConfiguration{RPID: "example.com", RPOrigin: "https://example.com"},
		SessionStore:   webAuthnSessionStore,
	}

	options, err := webAuthnSessionStore.BeginWebAuthnLogin(user, nil, sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(options)
}

func finishWebAuthnLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email        string
		WebAuthnData string
	}
	json.NewDecoder(r.Body).Decode(&req)

	user := User{Email: req.Email}
	sr := SessionRequest{
		WebAuthnConfig: WebAuthnConfiguration{RPID: "example.com", RPOrigin: "https://example.com"},
		SessionStore:   webAuthnSessionStore,
		WebAuthnData:   req.WebAuthnData,
	}

	err := webAuthnSessionStore.FinishWebAuthnLogin(user, nil, sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
