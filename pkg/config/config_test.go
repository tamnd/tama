package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

// clearEnv shields the test from TAMA_* values in the developer's shell.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"TAMA_ADDR", "TAMA_DATA", "TAMA_LLM_BASE_URL", "TAMA_LLM_API_KEY", "TAMA_LLM_MODEL", "TAMA_LOG_LEVEL"} {
		t.Setenv(k, "")
	}
}

// newFlags mirrors the flags serve registers, optionally pre-set.
func newFlags(t *testing.T, set map[string]string) *pflag.FlagSet {
	t.Helper()
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("addr", ":4321", "")
	fs.String("data", "", "")
	for k, v := range set {
		if err := fs.Set(k, v); err != nil {
			t.Fatal(err)
		}
	}
	return fs
}

func writeConfigFile(t *testing.T, dir, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func load(t *testing.T, flags *pflag.FlagSet) *Config {
	t.Helper()
	cfg, err := Load(flags)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return cfg
}

func TestDefaults(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	t.Setenv("TAMA_DATA", dir)

	cfg := load(t, nil)
	if cfg.Addr != ":4321" {
		t.Errorf("Addr = %q, want :4321", cfg.Addr)
	}
	if cfg.LLM.BaseURL != "http://127.0.0.1:8000/v1" {
		t.Errorf("BaseURL = %q", cfg.LLM.BaseURL)
	}
	if cfg.LLM.RequestTimeout != 5*time.Minute || cfg.LLM.ConnectTimeout != 30*time.Second {
		t.Errorf("timeouts = %v/%v", cfg.LLM.RequestTimeout, cfg.LLM.ConnectTimeout)
	}
	if cfg.Log.Level != "info" || cfg.Log.Format != "text" {
		t.Errorf("log = %+v", cfg.Log)
	}
	if got := cfg.Source("addr"); got != SourceDefault {
		t.Errorf("Source(addr) = %s, want default", got)
	}
}

func TestFlagBeatsEnv(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	t.Setenv("TAMA_DATA", dir)
	t.Setenv("TAMA_ADDR", ":7000")

	cfg := load(t, newFlags(t, map[string]string{"addr": ":8000"}))
	if cfg.Addr != ":8000" {
		t.Errorf("Addr = %q, want flag value :8000", cfg.Addr)
	}
	if got := cfg.Source("addr"); got != SourceFlag {
		t.Errorf("Source(addr) = %s, want flag", got)
	}
}

func TestEnvBeatsFile(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	t.Setenv("TAMA_DATA", dir)
	t.Setenv("TAMA_ADDR", ":7000")
	writeConfigFile(t, dir, "[server]\naddr = \":6000\"\n")

	cfg := load(t, nil)
	if cfg.Addr != ":7000" {
		t.Errorf("Addr = %q, want env value :7000", cfg.Addr)
	}
	if got := cfg.Source("addr"); got != SourceEnv {
		t.Errorf("Source(addr) = %s, want env", got)
	}
}

func TestFileBeatsDefault(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	t.Setenv("TAMA_DATA", dir)
	writeConfigFile(t, dir, `
[server]
addr = ":6000"

[llm]
model = "little-helper"
request_timeout = "2m"

[log]
level = "debug"
format = "json"
`)

	cfg := load(t, nil)
	if cfg.Addr != ":6000" {
		t.Errorf("Addr = %q, want file value :6000", cfg.Addr)
	}
	if cfg.LLM.Model != "little-helper" {
		t.Errorf("Model = %q", cfg.LLM.Model)
	}
	if cfg.LLM.RequestTimeout != 2*time.Minute {
		t.Errorf("RequestTimeout = %v", cfg.LLM.RequestTimeout)
	}
	if cfg.Log.Level != "debug" || cfg.Log.Format != "json" {
		t.Errorf("log = %+v", cfg.Log)
	}
	if got := cfg.Source("addr"); got != SourceFile {
		t.Errorf("Source(addr) = %s, want file", got)
	}
}

func TestLogLevelEnvJSONSuffix(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	t.Setenv("TAMA_DATA", dir)
	t.Setenv("TAMA_LOG_LEVEL", "warn,json")

	cfg := load(t, nil)
	if cfg.Log.Level != "warn" || cfg.Log.Format != "json" {
		t.Errorf("log = %+v, want warn/json", cfg.Log)
	}
}

func TestUnknownKeyWarns(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	t.Setenv("TAMA_DATA", dir)
	writeConfigFile(t, dir, "[server]\naddr = \":6000\"\nshiny = true\n")

	cfg := load(t, nil)
	if len(cfg.Warnings) != 1 || !strings.Contains(cfg.Warnings[0], "server.shiny") {
		t.Errorf("Warnings = %v, want one about server.shiny", cfg.Warnings)
	}
}

func TestValidateListsAllProblems(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	t.Setenv("TAMA_DATA", dir)
	t.Setenv("TAMA_ADDR", "no-port")
	t.Setenv("TAMA_LLM_BASE_URL", "not a url")
	t.Setenv("TAMA_LOG_LEVEL", "loud")

	_, err := Load(nil)
	if err == nil {
		t.Fatal("Load succeeded, want error")
	}
	for _, want := range []string{"addr", "base_url", "log.level"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q misses %q", err, want)
		}
	}
}

func TestPrintRedactsAPIKey(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	t.Setenv("TAMA_DATA", dir)
	t.Setenv("TAMA_LLM_API_KEY", "sk-very-secret")

	cfg := load(t, nil)
	var b strings.Builder
	cfg.Print(&b)
	out := b.String()
	if strings.Contains(out, "sk-very-secret") {
		t.Error("Print leaked the API key")
	}
	if !strings.Contains(out, `llm.api_key         = "***"  (env)`) {
		t.Errorf("Print output misses redacted key line:\n%s", out)
	}
}

func TestDataDirCreated(t *testing.T) {
	clearEnv(t)
	dir := filepath.Join(t.TempDir(), "nested", "tama")
	t.Setenv("TAMA_DATA", dir)

	load(t, nil)
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("data dir not created: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o700 {
		t.Errorf("data dir perm = %o, want 700", perm)
	}
}

func TestFileDataKeyMovesDataDir(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	moved := filepath.Join(dir, "moved")
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(filepath.Join(home, ".tama"), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	writeConfigFile(t, filepath.Join(home, ".tama"), "[server]\ndata = \""+moved+"\"\n")

	cfg := load(t, nil)
	if cfg.DataDir != moved {
		t.Errorf("DataDir = %q, want %q", cfg.DataDir, moved)
	}
	if got := cfg.Source("data"); got != SourceFile {
		t.Errorf("Source(data) = %s, want file", got)
	}
}
