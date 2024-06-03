package main

import (
	"github.com/jmoiron/sqlx/types"
)

// User represents a user in the system.
type User struct {
	Email       string
	DisplayName string
}

// WebAuthnConfiguration holds the configuration for WebAuthn operations.
type WebAuthnConfiguration struct {
	RPID     string
	RPOrigin string
}

// WebAuthn holds the credentials for a user.
type WebAuthn struct {
	Email         string
	PublicKeyData types.JSONText
}

// SessionRequest holds information needed for session handling.
type SessionRequest struct {
	WebAuthnConfig WebAuthnConfiguration
	SessionStore   *WebAuthnSessionStore
	WebAuthnData   string
}
