package creds

import (
	"fmt"
	"strings"

	"github.com/zalando/go-keyring"
)

// Service is the keyring service name, matching adif-mcp's SERVICE constant.
const Service = "adif-mcp"

// subject builds the keyring username: "{persona}:{provider}" lowercased.
func subject(persona, provider string) string {
	return strings.ToLower(persona + ":" + provider)
}

// Set stores credentials in the OS keyring.
func Set(persona, provider string, c *Credentials) error {
	data, err := c.ToJSON()
	if err != nil {
		return fmt.Errorf("serializing credentials: %w", err)
	}
	return keyring.Set(Service, subject(persona, provider), data)
}

// Get retrieves credentials from the OS keyring.
func Get(persona, provider string) (*Credentials, error) {
	raw, err := keyring.Get(Service, subject(persona, provider))
	if err != nil {
		return nil, err
	}
	return FromJSON(raw)
}

// Delete removes credentials from the OS keyring.
func Delete(persona, provider string) error {
	return keyring.Delete(Service, subject(persona, provider))
}

// Exists checks whether credentials are stored for a persona:provider pair.
func Exists(persona, provider string) bool {
	_, err := keyring.Get(Service, subject(persona, provider))
	return err == nil
}
