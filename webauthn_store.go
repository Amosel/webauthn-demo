package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pkg/errors"
)

type WebAuthnSessionStore struct {
	inProgressRegistrations map[string]string
	mu                      sync.Mutex
}

func NewWebAuthnSessionStore() *WebAuthnSessionStore {
	return &WebAuthnSessionStore{
		inProgressRegistrations: map[string]string{},
	}
}

func (store *WebAuthnSessionStore) SaveWebauthnSession(key string, data *webauthn.SessionData) error {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling session data: %v", err)
		return err
	}
	log.Printf("Session data saved for key: %s", key)
	store.put(key, string(marshaledData))
	return nil
}

func (store *WebAuthnSessionStore) put(key, val string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.inProgressRegistrations[key] = val
}

func (store *WebAuthnSessionStore) take(key string) (val string, ok bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	val, ok = store.inProgressRegistrations[key]
	if ok {
		delete(store.inProgressRegistrations, key)
	}
	return
}

func (store *WebAuthnSessionStore) GetWebauthnSession(key string) (data webauthn.SessionData, err error) {
	assertion, ok := store.take(key)
	if !ok {
		err = errors.New("assertion not in challenge store")
		return
	}
	err = json.Unmarshal([]byte(assertion), &data)
	return
}

func (store *WebAuthnSessionStore) BeginWebAuthnRegistration(user User, uwas []WebAuthn, config WebAuthnConfiguration) (*protocol.CredentialCreation, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Example RP",    // Display Name
		RPID:          config.RPID,     // Generally the domain name
		RPOrigin:      config.RPOrigin, // The origin URL for WebAuthn requests
	})

	if err != nil {
		return nil, err
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return nil, err
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = waUser.CredentialExcludeList()
	}

	options, sessionData, err := webAuthn.BeginRegistration(
		waUser,
		registerOptions,
	)

	if err != nil {
		return nil, err
	}

	userRegistrationIndexKey := fmt.Sprintf("%s-registration", user.Email)
	err = store.SaveWebauthnSession(userRegistrationIndexKey, sessionData)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (store *WebAuthnSessionStore) FinishWebAuthnRegistration(user User, uwas []WebAuthn, response *http.Request, config WebAuthnConfiguration) (*webauthn.Credential, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Example RP",    // Display Name
		RPID:          config.RPID,     // Generally the domain name
		RPOrigin:      config.RPOrigin, // The origin URL for WebAuthn requests
	})
	if err != nil {
		return nil, err
	}

	userRegistrationIndexKey := fmt.Sprintf("%s-registration", user.Email)
	sessionData, err := store.GetWebauthnSession(userRegistrationIndexKey)
	if err != nil {
		return nil, err
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return nil, err
	}

	credential, err := webAuthn.FinishRegistration(waUser, sessionData, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to FinishRegistration")
	}

	return credential, nil
}

func (store *WebAuthnSessionStore) BeginWebAuthnLogin(user User, uwas []WebAuthn, sr SessionRequest) (*protocol.CredentialAssertion, error) {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Example RP",               // Display Name
		RPID:          sr.WebAuthnConfig.RPID,     // Generally the domain name
		RPOrigin:      sr.WebAuthnConfig.RPOrigin, // The origin URL for WebAuthn requests
	})

	if err != nil {
		return nil, err
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return nil, err
	}

	options, sessionData, err := webAuthn.BeginLogin(waUser)
	if err != nil {
		return nil, err
	}

	userLoginIndexKey := fmt.Sprintf("%s-authentication", user.Email)
	err = sr.SessionStore.SaveWebauthnSession(userLoginIndexKey, sessionData)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (store *WebAuthnSessionStore) FinishWebAuthnLogin(user User, uwas []WebAuthn, sr SessionRequest) error {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Example RP",               // Display Name
		RPID:          sr.WebAuthnConfig.RPID,     // Generally the domain name
		RPOrigin:      sr.WebAuthnConfig.RPOrigin, // The origin URL for WebAuthn requests
	})

	if err != nil {
		return errors.Wrapf(err, "failed to create webAuthn structure with RPID: %s and RPOrigin: %s", sr.WebAuthnConfig.RPID, sr.WebAuthnConfig.RPOrigin)
	}

	credential, err := protocol.ParseCredentialRequestResponseBody(strings.NewReader(sr.WebAuthnData))
	if err != nil {
		return err
	}

	userLoginIndexKey := fmt.Sprintf("%s-authentication", user.Email)
	sessionData, err := sr.SessionStore.GetWebauthnSession(userLoginIndexKey)
	if err != nil {
		return err
	}

	waUser, err := duoWebAuthUserFromUser(user, uwas)
	if err != nil {
		return err
	}

	_, err = webAuthn.ValidateLogin(waUser, sessionData, credential)
	return err
}
