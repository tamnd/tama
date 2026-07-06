package server

import (
	"os/exec"
	"strings"
	"testing"
)

const module = "github.com/tamnd/tama"

// goListDeps returns the module-local deps of a package, itself excluded.
func goListDeps(t *testing.T, pkg string) []string {
	t.Helper()
	out, err := exec.Command("go", "list", "-deps", module+"/"+pkg).Output()
	if err != nil {
		t.Fatalf("go list -deps %s: %v", pkg, err)
	}
	var deps []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.HasPrefix(line, module+"/") && line != module+"/"+pkg {
			deps = append(deps, strings.TrimPrefix(line, module+"/"))
		}
	}
	return deps
}

// TestImportDirection pins the layering: store sits at the bottom, api may
// reach store and the domain packages, and nothing imports cmd.
func TestImportDirection(t *testing.T) {
	for _, dep := range goListDeps(t, "pkg/store") {
		if strings.HasPrefix(dep, "pkg/") {
			t.Errorf("pkg/store imports sibling %s; store must stay at the bottom", dep)
		}
	}

	apiAllowed := map[string]bool{
		"pkg/store":    true,
		"pkg/course":   true,
		"pkg/exercise": true,
		"pkg/engine":   true,
		"pkg/gen":      true,
	}
	for _, dep := range goListDeps(t, "pkg/api") {
		if !apiAllowed[dep] {
			t.Errorf("pkg/api imports %s, outside its allowed set", dep)
		}
	}

	for _, pkg := range []string{"pkg/api", "pkg/cli", "pkg/config", "pkg/course", "pkg/engine", "pkg/exercise", "pkg/gen", "pkg/server", "pkg/store", "web"} {
		for _, dep := range goListDeps(t, pkg) {
			if strings.HasPrefix(dep, "cmd/") {
				t.Errorf("%s imports %s; nothing imports cmd", pkg, dep)
			}
		}
	}
}
