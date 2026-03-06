package persona

// ProviderRef matches adif-mcp's ProviderRef TypedDict.
type ProviderRef struct {
	Username string `json:"username"`
}

// Persona matches adif-mcp's Persona TypedDict.
// Start/End use *string to preserve null in JSON (Python writes null, not omit).
type Persona struct {
	Name      string                 `json:"name"`
	Callsign  string                 `json:"callsign"`
	Start     *string                `json:"start"`
	End       *string                `json:"end"`
	Providers map[string]ProviderRef `json:"providers"`
}

// StartStr returns the start date as a string (empty if nil).
func (p *Persona) StartStr() string {
	if p.Start == nil {
		return ""
	}
	return *p.Start
}

// EndStr returns the end date as a string (empty if nil).
func (p *Persona) EndStr() string {
	if p.End == nil {
		return ""
	}
	return *p.End
}

// StringPtr returns a pointer to s, or nil if empty.
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// PersonaIndex is the top-level JSON structure of personas.json.
type PersonaIndex struct {
	Personas map[string]Persona `json:"personas"`
}

// PersonaSummary is a frontend-friendly view of a persona.
type PersonaSummary struct {
	Name      string   `json:"name"`
	Callsign  string   `json:"callsign"`
	Start     string   `json:"start"`
	End       string   `json:"end"`
	Providers []string `json:"providers"`
	IsActive  bool     `json:"isActive"`
}
