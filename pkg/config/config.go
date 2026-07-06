// Package config resolves tama's runtime settings. Precedence is flags over
// TAMA_* environment variables over defaults; flag overrides are applied by
// the cli layer after Load.
package config

import (
	"os"
	"path/filepath"
)

// Config holds everything the server and the generator need to run.
type Config struct {
	// Addr is the HTTP listen address, ":4321" by default.
	Addr string
	// DataDir holds the SQLite database and the course packs, ~/.tama by
	// default. It is created on first use.
	DataDir string
	// LLM points at the OpenAI-compatible endpoint used for course
	// generation. The lesson path never touches it.
	LLM LLM
}

// LLM is the connection to the OpenAI-compatible generation endpoint.
type LLM struct {
	BaseURL string
	APIKey  string
	Model   string
}

// Load builds a Config from environment variables and defaults.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Addr:    envOr("TAMA_ADDR", ":4321"),
		DataDir: envOr("TAMA_DATA", filepath.Join(home, ".tama")),
		LLM: LLM{
			BaseURL: envOr("TAMA_LLM_BASE_URL", "http://127.0.0.1:8000/v1"),
			APIKey:  os.Getenv("TAMA_LLM_API_KEY"),
			Model:   os.Getenv("TAMA_LLM_MODEL"),
		},
	}
	return cfg, nil
}

// DBPath is the SQLite database file inside the data directory.
func (c *Config) DBPath() string {
	return filepath.Join(c.DataDir, "tama.db")
}

// PacksDir is where generated course packs are cached.
func (c *Config) PacksDir() string {
	return filepath.Join(c.DataDir, "packs")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
