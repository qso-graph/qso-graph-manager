package persona

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Load reads and parses personas.json. Returns empty index if file doesn't exist.
func Load() (*PersonaIndex, error) {
	path := PersonasPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &PersonaIndex{Personas: make(map[string]Persona)}, nil
		}
		return nil, fmt.Errorf("reading personas.json: %w", err)
	}
	var idx PersonaIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parsing personas.json: %w", err)
	}
	if idx.Personas == nil {
		idx.Personas = make(map[string]Persona)
	}
	// Normalize map keys to lowercase — Python may write uppercase keys
	for k, v := range idx.Personas {
		lower := strings.ToLower(k)
		if lower != k {
			delete(idx.Personas, k)
			idx.Personas[lower] = v
		}
	}
	return &idx, nil
}

// Save writes the persona index to personas.json with sorted keys and indent=2,
// matching adif-mcp's Python output.
func Save(idx *PersonaIndex) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing personas: %w", err)
	}
	data = append(data, '\n')
	return os.WriteFile(PersonasPath(), data, 0o644)
}

// Add creates a new persona in the index.
func Add(idx *PersonaIndex, name, callsign, start, end string) error {
	key := strings.ToLower(name)
	if _, exists := idx.Personas[key]; exists {
		return fmt.Errorf("persona %q already exists", name)
	}
	idx.Personas[key] = Persona{
		Name:      name,
		Callsign:  strings.ToUpper(callsign),
		Start:     StringPtr(start),
		End:       StringPtr(end),
		Providers: make(map[string]ProviderRef),
	}
	return Save(idx)
}

// Remove deletes a persona from the index.
func Remove(idx *PersonaIndex, name string) error {
	key := strings.ToLower(name)
	if _, exists := idx.Personas[key]; !exists {
		return fmt.Errorf("persona %q not found", name)
	}
	delete(idx.Personas, key)
	return Save(idx)
}

// SetActive writes the persona name to current.txt.
func SetActive(name string) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	return os.WriteFile(CurrentPath(), []byte(name), 0o644)
}

// GetActive reads the active persona name from current.txt.
func GetActive() string {
	data, err := os.ReadFile(CurrentPath())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// EnableProvider adds a provider reference to a persona.
func EnableProvider(idx *PersonaIndex, personaName, provider, username string) error {
	key := strings.ToLower(personaName)
	p, exists := idx.Personas[key]
	if !exists {
		return fmt.Errorf("persona %q not found", personaName)
	}
	if p.Providers == nil {
		p.Providers = make(map[string]ProviderRef)
	}
	p.Providers[strings.ToLower(provider)] = ProviderRef{Username: username}
	idx.Personas[key] = p
	return Save(idx)
}

// DisableProvider removes a provider reference from a persona.
func DisableProvider(idx *PersonaIndex, personaName, provider string) error {
	key := strings.ToLower(personaName)
	p, exists := idx.Personas[key]
	if !exists {
		return fmt.Errorf("persona %q not found", personaName)
	}
	delete(p.Providers, strings.ToLower(provider))
	idx.Personas[key] = p
	return Save(idx)
}

// Summaries returns a list of PersonaSummary for the frontend.
func Summaries(idx *PersonaIndex) []PersonaSummary {
	active := GetActive()
	out := make([]PersonaSummary, 0, len(idx.Personas))
	for _, p := range idx.Personas {
		providers := make([]string, 0, len(p.Providers))
		for k := range p.Providers {
			providers = append(providers, k)
		}
		out = append(out, PersonaSummary{
			Name:      p.Name,
			Callsign:  p.Callsign,
			Start:     p.StartStr(),
			End:       p.EndStr(),
			Providers: providers,
			IsActive:  strings.EqualFold(p.Name, active),
		})
	}
	return out
}

// GetByName retrieves a persona by name (case-insensitive).
func GetByName(idx *PersonaIndex, name string) (*Persona, error) {
	key := strings.ToLower(name)
	p, exists := idx.Personas[key]
	if !exists {
		return nil, fmt.Errorf("persona %q not found", name)
	}
	return &p, nil
}

// ProviderUsername returns the username for a provider in a persona, or empty string.
func ProviderUsername(idx *PersonaIndex, personaName, provider string) string {
	key := strings.ToLower(personaName)
	p, exists := idx.Personas[key]
	if !exists {
		return ""
	}
	ref, exists := p.Providers[strings.ToLower(provider)]
	if !exists {
		return ""
	}
	return ref.Username
}

// ListProviderKeys returns all provider slugs for a persona.
func ListProviderKeys(p *Persona) []string {
	keys := make([]string, 0, len(p.Providers))
	for k := range p.Providers {
		keys = append(keys, k)
	}
	return keys
}

// HasProvider returns the directory for a given persona and provider.
func HasProvider(idx *PersonaIndex, personaName, provider string) bool {
	key := strings.ToLower(personaName)
	p, exists := idx.Personas[key]
	if !exists {
		return false
	}
	provKey := strings.ToLower(provider)
	if provKey == "qrz_logbook" {
		provKey = "qrz_logbook"
	}
	_, has := p.Providers[provKey]
	return has
}

// EnsureConfigDir creates the config directory if it doesn't exist.
func EnsureConfigDir() error {
	return os.MkdirAll(ConfigDir(), 0o755)
}

// ConfigDirPath returns the config dir path (exported for frontend).
func ConfigDirPath() string {
	return filepath.Clean(ConfigDir())
}
