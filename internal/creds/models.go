package creds

import "encoding/json"

// Credentials matches adif-mcp's Credentials dataclass exactly.
type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	APIKey   string `json:"api_key,omitempty"`
}

// ToJSON serializes credentials to JSON matching adif-mcp's format.
func (c *Credentials) ToJSON() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON deserializes credentials from JSON.
func FromJSON(raw string) (*Credentials, error) {
	var c Credentials
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// CredentialInfo is a frontend-safe view (secret masked).
type CredentialInfo struct {
	Persona  string `json:"persona"`
	Provider string `json:"provider"`
	Username string `json:"username"`
	HasCreds bool   `json:"hasCreds"`
	AuthType string `json:"authType"`
}

// HealthEntry represents one persona×provider health check result.
type HealthEntry struct {
	Persona  string `json:"persona"`
	Provider string `json:"provider"`
	Username string `json:"username"`
	HasCreds bool   `json:"hasCreds"`
	AuthType string `json:"authType"`
}
