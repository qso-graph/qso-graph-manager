package persona

import (
	"os"
	"path/filepath"
	"runtime"
)

// ConfigDir returns the adif-mcp config directory, matching Python's paths.py logic.
func ConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, "adif-mcp")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "AppData", "Roaming", "adif-mcp")
	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, "adif-mcp")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", "adif-mcp")
	}
}

// PersonasPath returns the full path to personas.json.
func PersonasPath() string {
	return filepath.Join(ConfigDir(), "personas.json")
}

// CurrentPath returns the full path to current.txt (active persona marker).
func CurrentPath() string {
	return filepath.Join(ConfigDir(), "current.txt")
}
