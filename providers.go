package main

// ProviderInfo describes a qso-graph MCP server's auth requirements.
type ProviderInfo struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	AuthType string `json:"authType"` // "username_password", "api_key", or "none"
	PyPI     string `json:"pypi"`
}

// KnownProviders is the embedded provider metadata for all qso-graph MCP servers.
var KnownProviders = []ProviderInfo{
	{Name: "eQSL", Slug: "eqsl", AuthType: "username_password", PyPI: "eqsl-mcp"},
	{Name: "QRZ XML", Slug: "qrz", AuthType: "username_password", PyPI: "qrz-mcp"},
	{Name: "QRZ Logbook", Slug: "qrz_logbook", AuthType: "api_key", PyPI: "qrz-mcp"},
	{Name: "LoTW", Slug: "lotw", AuthType: "username_password", PyPI: "lotw-mcp"},
	{Name: "HamQTH", Slug: "hamqth", AuthType: "username_password", PyPI: "hamqth-mcp"},
}

// providerBySlug looks up a provider by slug.
func providerBySlug(slug string) *ProviderInfo {
	for i := range KnownProviders {
		if KnownProviders[i].Slug == slug {
			return &KnownProviders[i]
		}
	}
	return nil
}
