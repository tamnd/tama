// Package config resolves tama's runtime settings from four layers, highest
// wins: command-line flags, TAMA_* environment variables, config.toml in the
// data directory, and built-in defaults. Nothing else in the tree reads
// os.Getenv; every consumer gets a resolved Config from Load.
package config

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"
)

// Source says which layer a resolved value came from, for --print-config.
type Source string

// The four layers, in precedence order.
const (
	SourceFlag    Source = "flag"
	SourceEnv     Source = "env"
	SourceFile    Source = "file"
	SourceDefault Source = "default"
)

// Config holds everything the server and the generator need to run.
type Config struct {
	// Addr is the HTTP listen address, ":4321" by default.
	Addr string
	// DataDir holds the SQLite database, config.toml, and generated audio,
	// ~/.tama by default. Load creates it 0700 if missing.
	DataDir string
	// LLM points at the OpenAI-compatible endpoint used for course
	// generation. The lesson path never touches it.
	LLM LLM
	// Log controls slog level and output format.
	Log Log
	// Warnings collects non-fatal problems found while loading, like unknown
	// keys in config.toml. The caller decides where to surface them.
	Warnings []string

	sources map[string]Source
}

// LLM is the connection to the OpenAI-compatible generation endpoint.
type LLM struct {
	BaseURL string
	APIKey  string
	Model   string
	// RequestTimeout bounds one whole request, 5m by default.
	RequestTimeout time.Duration
	// ConnectTimeout bounds connect-to-first-byte, 30s by default.
	ConnectTimeout time.Duration
}

// Log holds the logging knobs.
type Log struct {
	Level  string // debug, info, warn, error
	Format string // text or json
}

// fileConfig mirrors config.toml. Pointers tell a set key from a missing one.
type fileConfig struct {
	Server struct {
		Addr *string `toml:"addr"`
		Data *string `toml:"data"`
	} `toml:"server"`
	LLM struct {
		BaseURL        *string `toml:"base_url"`
		APIKey         *string `toml:"api_key"`
		Model          *string `toml:"model"`
		RequestTimeout *string `toml:"request_timeout"`
		ConnectTimeout *string `toml:"connect_timeout"`
	} `toml:"llm"`
	Log struct {
		Level  *string `toml:"level"`
		Format *string `toml:"format"`
	} `toml:"log"`
}

// Load resolves the config. flags may be nil for commands that define none;
// when present, only flags the user actually set override the lower layers.
// Validation reports every problem at once rather than the first hit.
func Load(flags *pflag.FlagSet) (*Config, error) {
	cfg := &Config{
		Addr: ":4321",
		LLM: LLM{
			BaseURL:        "http://127.0.0.1:8000/v1",
			RequestTimeout: 5 * time.Minute,
			ConnectTimeout: 30 * time.Second,
		},
		Log:     Log{Level: "info", Format: "text"},
		sources: map[string]Source{},
	}
	for _, k := range keys {
		cfg.sources[k] = SourceDefault
	}

	// The data dir hosts config.toml, so it resolves first, from flag over
	// env over default. A data key inside the file still wins over the
	// default for everything that runs after Load.
	dataDir, dataSource, err := resolveDataDir(flags)
	if err != nil {
		return nil, err
	}
	cfg.DataDir = dataDir
	cfg.sources["data"] = dataSource

	if err := cfg.applyFile(); err != nil {
		return nil, err
	}
	cfg.applyEnv()
	cfg.applyFlags(flags)

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// keys names every resolved value, in the order Print shows them.
var keys = []string{
	"addr", "data",
	"llm.base_url", "llm.api_key", "llm.model",
	"llm.request_timeout", "llm.connect_timeout",
	"log.level", "log.format",
}

func resolveDataDir(flags *pflag.FlagSet) (string, Source, error) {
	if f := changedFlag(flags, "data"); f != "" {
		return expandHome(f), SourceFlag, nil
	}
	if v, ok := os.LookupEnv("TAMA_DATA"); ok && v != "" {
		return expandHome(v), SourceEnv, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", SourceDefault, err
	}
	return filepath.Join(home, ".tama"), SourceDefault, nil
}

// applyFile layers config.toml over the defaults. A missing file is fine; a
// malformed one is an error; an unknown key is only a warning so files written
// for a newer tama still load on an older binary.
func (c *Config) applyFile() error {
	path := c.ConfigPath()
	var fc fileConfig
	md, err := toml.DecodeFile(path, &fc)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	for _, k := range md.Undecoded() {
		c.Warnings = append(c.Warnings, fmt.Sprintf("%s: unknown key %q ignored", path, k.String()))
	}

	set := func(key string, dst *string, v *string) {
		if v != nil {
			*dst = *v
			c.sources[key] = SourceFile
		}
	}
	set("addr", &c.Addr, fc.Server.Addr)
	if fc.Server.Data != nil && c.sources["data"] == SourceDefault {
		c.DataDir = expandHome(*fc.Server.Data)
		c.sources["data"] = SourceFile
	}
	set("llm.base_url", &c.LLM.BaseURL, fc.LLM.BaseURL)
	set("llm.api_key", &c.LLM.APIKey, fc.LLM.APIKey)
	set("llm.model", &c.LLM.Model, fc.LLM.Model)
	set("log.level", &c.Log.Level, fc.Log.Level)
	set("log.format", &c.Log.Format, fc.Log.Format)

	if err := c.setDuration("llm.request_timeout", &c.LLM.RequestTimeout, fc.LLM.RequestTimeout, path); err != nil {
		return err
	}
	return c.setDuration("llm.connect_timeout", &c.LLM.ConnectTimeout, fc.LLM.ConnectTimeout, path)
}

func (c *Config) setDuration(key string, dst *time.Duration, v *string, path string) error {
	if v == nil {
		return nil
	}
	d, err := time.ParseDuration(*v)
	if err != nil {
		return fmt.Errorf("%s: %s: %w", path, key, err)
	}
	*dst = d
	c.sources[key] = SourceFile
	return nil
}

func (c *Config) applyEnv() {
	set := func(key, env string, dst *string) {
		if v, ok := os.LookupEnv(env); ok && v != "" {
			*dst = v
			c.sources[key] = SourceEnv
		}
	}
	set("addr", "TAMA_ADDR", &c.Addr)
	set("llm.base_url", "TAMA_LLM_BASE_URL", &c.LLM.BaseURL)
	set("llm.api_key", "TAMA_LLM_API_KEY", &c.LLM.APIKey)
	set("llm.model", "TAMA_LLM_MODEL", &c.LLM.Model)

	// TAMA_LOG_LEVEL is "debug" or "debug,json"; the suffix flips the
	// handler to JSON without a second variable.
	if v, ok := os.LookupEnv("TAMA_LOG_LEVEL"); ok && v != "" {
		level, isJSON := strings.CutSuffix(v, ",json")
		c.Log.Level = level
		c.sources["log.level"] = SourceEnv
		if isJSON {
			c.Log.Format = "json"
			c.sources["log.format"] = SourceEnv
		}
	}
}

func (c *Config) applyFlags(flags *pflag.FlagSet) {
	if f := changedFlag(flags, "addr"); f != "" {
		c.Addr = f
		c.sources["addr"] = SourceFlag
	}
	// data was already taken from the flag in resolveDataDir.
}

// validate collects every problem before failing, so a broken config file is
// one round trip to fix, not four.
func (c *Config) validate() error {
	var problems []string

	if _, _, err := net.SplitHostPort(c.Addr); err != nil {
		problems = append(problems, fmt.Sprintf("addr %q is not host:port: %v", c.Addr, err))
	}
	if err := os.MkdirAll(c.DataDir, 0o700); err != nil {
		problems = append(problems, fmt.Sprintf("data dir %s: %v", c.DataDir, err))
	} else if err := writable(c.DataDir); err != nil {
		problems = append(problems, fmt.Sprintf("data dir %s is not writable: %v", c.DataDir, err))
	}
	if u, err := url.Parse(c.LLM.BaseURL); err != nil {
		problems = append(problems, fmt.Sprintf("llm.base_url %q: %v", c.LLM.BaseURL, err))
	} else if u.Scheme != "http" && u.Scheme != "https" || u.Host == "" {
		problems = append(problems, fmt.Sprintf("llm.base_url %q must be an http(s) URL", c.LLM.BaseURL))
	}
	switch c.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		problems = append(problems, fmt.Sprintf("log.level %q must be debug, info, warn, or error", c.Log.Level))
	}
	switch c.Log.Format {
	case "text", "json":
	default:
		problems = append(problems, fmt.Sprintf("log.format %q must be text or json", c.Log.Format))
	}
	if c.LLM.RequestTimeout <= 0 {
		problems = append(problems, "llm.request_timeout must be positive")
	}
	if c.LLM.ConnectTimeout <= 0 {
		problems = append(problems, "llm.connect_timeout must be positive")
	}

	if len(problems) > 0 {
		return fmt.Errorf("config:\n  %s", strings.Join(problems, "\n  "))
	}
	return nil
}

// writable probes the directory with a throwaway file, catching read-only
// mounts that MkdirAll on an existing dir would not.
func writable(dir string) error {
	f, err := os.CreateTemp(dir, ".tama-write-*")
	if err != nil {
		return err
	}
	f.Close()
	return os.Remove(f.Name())
}

// Source reports which layer set the named value ("addr", "llm.model", ...).
func (c *Config) Source(key string) Source {
	return c.sources[key]
}

// Print writes the resolved config with per-value provenance, the body of
// `tama serve --print-config`. The API key never renders in the clear.
func (c *Config) Print(w io.Writer) {
	vals := map[string]string{
		"addr":                c.Addr,
		"data":                c.DataDir,
		"llm.base_url":        c.LLM.BaseURL,
		"llm.api_key":         c.LLM.APIKey,
		"llm.model":           c.LLM.Model,
		"llm.request_timeout": c.LLM.RequestTimeout.String(),
		"llm.connect_timeout": c.LLM.ConnectTimeout.String(),
		"log.level":           c.Log.Level,
		"log.format":          c.Log.Format,
	}
	for _, k := range keys {
		v := vals[k]
		if k == "llm.api_key" && v != "" {
			v = "***"
		}
		fmt.Fprintf(w, "%-19s = %q  (%s)\n", k, v, c.sources[k])
	}
}

// DBPath is the SQLite database file inside the data directory.
func (c *Config) DBPath() string {
	return filepath.Join(c.DataDir, "tama.db")
}

// ConfigPath is the optional TOML config file inside the data directory.
func (c *Config) ConfigPath() string {
	return filepath.Join(c.DataDir, "config.toml")
}

// AudioDir is where M8 caches generated audio; created lazily, not here.
func (c *Config) AudioDir() string {
	return filepath.Join(c.DataDir, "audio")
}

func changedFlag(flags *pflag.FlagSet, name string) string {
	if flags == nil {
		return ""
	}
	f := flags.Lookup(name)
	if f == nil || !f.Changed {
		return ""
	}
	return f.Value.String()
}

func expandHome(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[1:])
		}
	}
	return p
}
