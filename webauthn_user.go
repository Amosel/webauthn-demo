package main

import (
	"fmt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthnUser implements the required WebAuthn user interface
type WebAuthnUser struct {
	Email         string
	WACredentials []webauthn.Credential
}

// WebAuthnID returns the user's ID
func (u WebAuthnUser) WebAuthnID() []byte {
	return []byte(u.Email)
}

// WebAuthnName returns the user's email
func (u WebAuthnUser) WebAuthnName() string {
	return u.Email
}

// WebAuthnDisplayName returns the user's display name.
func (u WebAuthnUser) WebAuthnDisplayName() string {
	return u.Email
}

// WebAuthnIcon returns the user's icon.
func (u WebAuthnUser) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials returns credentials owned by the user
func (u WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.WACredentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials to prevent them from re-registering keys
func (u WebAuthnUser) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.WACredentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}
	return credentialExcludeList
}

func (u *WebAuthnUser) LoadWebAuthnCredentials(uwas []WebAuthn) error {
	for _, v := range uwas {
		var credential webauthn.Credential
		err := v.PublicKeyData.Unmarshal(&credential)
		if err != nil {
			return fmt.Errorf("error unmarshalling provided PublicKeyData: %s", err)
		}
		u.WACredentials = append(u.WACredentials, credential)
	}
	return nil
}

func duoWebAuthUserFromUser(user User, uwas []WebAuthn) (WebAuthnUser, error) {
	waUser := WebAuthnUser{
		Email: user.Email,
	}
	err := waUser.LoadWebAuthnCredentials(uwas)

	return waUser, err
}
