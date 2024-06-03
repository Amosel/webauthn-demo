package main

import (
	"regexp"
)

// isValidEmail checks if the provided string is a valid email address.
func isValidEmail(email string) bool {
	// Regular expression for validating an email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
