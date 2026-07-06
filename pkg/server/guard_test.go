package server

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// repoRoot climbs from the test's working directory to the module root.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the test directory")
		}
		dir = parent
	}
}

// TestNoInternalDirectories fails the build if any internal/ directory
// appears anywhere in the repo; the layout keeps everything under pkg/.
func TestNoInternalDirectories(t *testing.T) {
	root := repoRoot(t)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return err
		}
		switch d.Name() {
		case ".git", "node_modules":
			return filepath.SkipDir
		case "internal":
			t.Errorf("internal directory found at %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// depAllowlist is the reviewed set of direct dependencies. Growing it is a
// deliberate diff here, next to a row in CONTRIBUTING.md.
var depAllowlist = map[string]bool{
	"modernc.org/sqlite":            true,
	"github.com/spf13/cobra":        true,
	"github.com/spf13/pflag":        true,
	"github.com/charmbracelet/fang": true,
	"golang.org/x/crypto":           true,
	"github.com/BurntSushi/toml":    true,
	"github.com/klauspost/compress": true,
}

// TestDirectDependenciesAreAllowlisted parses go.mod and fails on any direct
// require outside the allowlist. Indirect modules (fang's charm stack and
// sqlite's build deps) are fine.
func TestDirectDependenciesAreAllowlisted(t *testing.T) {
	f, err := os.Open(filepath.Join(repoRoot(t), "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	inRequire := false
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "require ("):
			inRequire = true
			continue
		case line == ")":
			inRequire = false
			continue
		}
		if !inRequire || line == "" || strings.HasSuffix(line, "// indirect") {
			continue
		}
		mod, _, ok := strings.Cut(line, " ")
		if !ok {
			continue
		}
		if !depAllowlist[mod] {
			t.Errorf("direct dependency %s is not in the allowlist; add a row to CONTRIBUTING.md and this test", mod)
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
}

// slogCall matches the slog call sites whose key/value pairs we screen.
var slogCall = regexp.MustCompile(`slog\.\w+|\.(Info|Warn|Error|Debug|Log)\w*\(`)

// forbiddenLogField spots password or token used as a log field name.
var forbiddenLogField = regexp.MustCompile(`"(password|token|password_hash|session_token)"`)

// TestNoSecretLogFields greps every non-test Go file for log calls carrying
// a password or token field. Usernames may log; secrets never do.
func TestNoSecretLogFields(t *testing.T) {
	root := repoRoot(t)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == "web" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(body), "\n") {
			if slogCall.MatchString(line) && forbiddenLogField.MatchString(line) {
				t.Errorf("%s:%d logs a secret field: %s", path, i+1, strings.TrimSpace(line))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
