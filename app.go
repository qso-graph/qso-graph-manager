package main

import (
	"context"
	"fmt"
	"strings"

	"qso-graph-manager/internal/creds"
	"qso-graph-manager/internal/persona"
)

// App is the main application struct bound to the Wails frontend.
type App struct {
	ctx context.Context
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{}
}

// startup is called at application start.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// --- Persona Management ---

// ListPersonas returns all personas with health summary info.
func (a *App) ListPersonas() ([]persona.PersonaSummary, error) {
	idx, err := persona.Load()
	if err != nil {
		return nil, err
	}
	return persona.Summaries(idx), nil
}

// GetPersona returns full details for a single persona.
func (a *App) GetPersona(name string) (*persona.Persona, error) {
	idx, err := persona.Load()
	if err != nil {
		return nil, err
	}
	return persona.GetByName(idx, name)
}

// AddPersona creates a new persona.
func (a *App) AddPersona(name, callsign, start, end string) error {
	if name == "" || callsign == "" {
		return fmt.Errorf("name and callsign are required")
	}
	idx, err := persona.Load()
	if err != nil {
		return err
	}
	return persona.Add(idx, name, callsign, start, end)
}

// RemovePersona deletes a persona and all its keyring credentials.
func (a *App) RemovePersona(name string) error {
	idx, err := persona.Load()
	if err != nil {
		return err
	}
	p, err := persona.GetByName(idx, name)
	if err != nil {
		return err
	}
	// Clean up keyring entries for all providers
	for provider := range p.Providers {
		_ = creds.Delete(name, provider) // best-effort
	}
	return persona.Remove(idx, name)
}

// SetActivePersona sets the active persona (writes current.txt).
func (a *App) SetActivePersona(name string) error {
	idx, err := persona.Load()
	if err != nil {
		return err
	}
	if _, err := persona.GetByName(idx, name); err != nil {
		return err
	}
	return persona.SetActive(name)
}

// GetActivePersona returns the currently active persona name.
func (a *App) GetActivePersona() string {
	return persona.GetActive()
}

// --- Credential Management ---

// SetCredential stores credentials in the OS keyring for a persona+provider.
func (a *App) SetCredential(personaName, provider, username, secret string) error {
	if personaName == "" || provider == "" {
		return fmt.Errorf("persona and provider are required")
	}

	// Determine auth type from provider metadata
	prov := providerBySlug(provider)
	authType := "username_password"
	if prov != nil {
		authType = prov.AuthType
	}

	var c creds.Credentials
	c.Username = username
	if authType == "api_key" {
		c.APIKey = secret
	} else {
		c.Password = secret
	}

	if err := creds.Set(personaName, provider, &c); err != nil {
		return fmt.Errorf("storing credentials: %w", err)
	}

	// Ensure provider is enabled in persona
	idx, err := persona.Load()
	if err != nil {
		return err
	}
	if !persona.HasProvider(idx, personaName, provider) {
		return persona.EnableProvider(idx, personaName, provider, username)
	}
	return nil
}

// GetCredential returns masked credential info (no secrets exposed to frontend).
func (a *App) GetCredential(personaName, provider string) (*creds.CredentialInfo, error) {
	prov := providerBySlug(provider)
	authType := "username_password"
	if prov != nil {
		authType = prov.AuthType
	}

	c, err := creds.Get(personaName, provider)
	if err != nil {
		return &creds.CredentialInfo{
			Persona:  personaName,
			Provider: provider,
			HasCreds: false,
			AuthType: authType,
		}, nil
	}

	return &creds.CredentialInfo{
		Persona:  personaName,
		Provider: provider,
		Username: c.Username,
		HasCreds: true,
		AuthType: authType,
	}, nil
}

// DeleteCredential removes credentials from the OS keyring.
func (a *App) DeleteCredential(personaName, provider string) error {
	return creds.Delete(personaName, provider)
}

// DoctorCheck runs a health check across all personas and providers.
func (a *App) DoctorCheck() ([]creds.HealthEntry, error) {
	idx, err := persona.Load()
	if err != nil {
		return nil, err
	}

	var entries []creds.HealthEntry
	for _, p := range idx.Personas {
		for _, prov := range KnownProviders {
			if prov.AuthType == "none" {
				continue
			}
			has := creds.Exists(p.Name, prov.Slug)
			username := ""
			if ref, ok := p.Providers[prov.Slug]; ok {
				username = ref.Username
			}
			entries = append(entries, creds.HealthEntry{
				Persona:  p.Name,
				Provider: prov.Slug,
				Username: username,
				HasCreds: has,
				AuthType: prov.AuthType,
			})
		}
	}
	return entries, nil
}

// --- Provider Management ---

// ListProviders returns all known providers with their auth types.
func (a *App) ListProviders() []ProviderInfo {
	return KnownProviders
}

// EnableProvider adds a provider to a persona.
func (a *App) EnableProvider(personaName, provider, username string) error {
	idx, err := persona.Load()
	if err != nil {
		return err
	}
	return persona.EnableProvider(idx, personaName, strings.ToLower(provider), username)
}

// DisableProvider removes a provider from a persona.
func (a *App) DisableProvider(personaName, provider string) error {
	idx, err := persona.Load()
	if err != nil {
		return err
	}
	return persona.DisableProvider(idx, personaName, strings.ToLower(provider))
}

// GetConfigDir returns the config directory path (for display in the UI).
func (a *App) GetConfigDir() string {
	return persona.ConfigDirPath()
}
