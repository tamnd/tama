package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/tamnd/tama/pkg/store"
)

// runCmd executes the root command with args and returns stdout.
func runCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := newRoot()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(args)
	err := root.ExecuteContext(context.Background())
	return out.String(), err
}

// TestEveryCommandHelp catches wiring regressions cheaply: --help must exit
// clean for the whole command surface.
func TestEveryCommandHelp(t *testing.T) {
	commands := [][]string{
		{"--help"},
		{"serve", "--help"},
		{"user", "--help"},
		{"user", "add", "--help"},
		{"user", "ls", "--help"},
		{"user", "rm", "--help"},
		{"user", "passwd", "--help"},
		{"db", "--help"},
		{"db", "migrate", "--help"},
		{"db", "status", "--help"},
		{"seed", "--help"},
		{"gen", "--help"},
		{"pack", "--help"},
		{"pack", "inspect", "--help"},
		{"pack", "validate", "--help"},
	}
	for _, args := range commands {
		out, err := runCmd(t, args...)
		if err != nil {
			t.Errorf("%v: %v", args, err)
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("%v: no usage in output", args)
		}
	}
}

// TestHelpHasExamples enforces the Short/Long/Example bar on every command.
func TestHelpHasExamples(t *testing.T) {
	for _, cmd := range newRoot().Commands() {
		if cmd.Short == "" || cmd.Long == "" {
			t.Errorf("%s: missing Short or Long", cmd.Name())
		}
		if len(cmd.Commands()) == 0 && cmd.Example == "" {
			t.Errorf("%s: missing Example", cmd.Name())
		}
		for _, sub := range cmd.Commands() {
			if sub.Short == "" || sub.Long == "" || sub.Example == "" {
				t.Errorf("%s %s: missing Short, Long, or Example", cmd.Name(), sub.Name())
			}
		}
	}
}

func TestUserCommands(t *testing.T) {
	data := t.TempDir()

	out, err := runCmd(t, "user", "add", "--data", data, "--username", "Mochi", "--password", "password1")
	if err != nil {
		t.Fatalf("user add: %v", err)
	}
	if !strings.Contains(out, "created user mochi") {
		t.Errorf("add output = %q", out)
	}

	out, err = runCmd(t, "user", "ls", "--data", data)
	if err != nil || !strings.Contains(out, "mochi") {
		t.Errorf("user ls = %q, %v", out, err)
	}

	if _, err = runCmd(t, "user", "passwd", "mochi", "--data", data, "--password", "password2"); err != nil {
		t.Errorf("user passwd: %v", err)
	}
	if _, err = runCmd(t, "user", "passwd", "ghost", "--data", data, "--password", "password2"); err == nil {
		t.Error("passwd for missing user succeeded")
	}

	if _, err = runCmd(t, "user", "rm", "mochi", "--data", data); err != nil {
		t.Errorf("user rm: %v", err)
	}
	if _, err = runCmd(t, "user", "rm", "mochi", "--data", data); err == nil {
		t.Error("second rm succeeded")
	}
}

func TestDBCommands(t *testing.T) {
	data := t.TempDir()

	out, err := runCmd(t, "db", "migrate", "--data", data)
	if err != nil || !strings.Contains(out, "applied 0001_core") {
		t.Errorf("db migrate = %q, %v", out, err)
	}
	out, err = runCmd(t, "db", "migrate", "--data", data)
	if err != nil || !strings.Contains(out, "up to date") {
		t.Errorf("second db migrate = %q, %v", out, err)
	}
	out, err = runCmd(t, "db", "status", "--data", data)
	if err != nil || !strings.Contains(out, "schema version 1, 0 pending") {
		t.Errorf("db status = %q, %v", out, err)
	}
}

func TestSeedDemoIsIdempotent(t *testing.T) {
	data := t.TempDir()
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		if _, err := runCmd(t, "seed", "--demo", "--data", data); err != nil {
			t.Fatalf("seed run %d: %v", i+1, err)
		}
	}

	db, err := store.Open(ctx, filepath.Join(data, "tama.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	user, err := db.UserByUsername(ctx, "demo")
	if err != nil {
		t.Fatalf("demo user: %v", err)
	}
	course, err := db.CourseByID(ctx, "es-en")
	if err != nil || course.Status != "ready" || course.PackID == nil {
		t.Errorf("course = %+v, %v", course, err)
	}
	pack, err := db.LatestPack(ctx, "es-en")
	if err != nil || pack.Version != 1 {
		t.Fatalf("pack = %+v, %v", pack, err)
	}

	// The stored blob is zstd; it must decompress back to the fixture.
	dec, _ := zstd.NewReader(nil)
	defer dec.Close()
	raw, err := dec.DecodeAll(pack.Content, nil)
	if err != nil || !bytes.Equal(raw, store.DemoPack()) {
		t.Errorf("pack content mismatch: %v", err)
	}

	total, err := db.XPTotal(ctx, user.ID)
	if err != nil || total != 10 {
		t.Errorf("xp after two seeds = %d, %v, want 10", total, err)
	}
	streak, err := db.GetStreak(ctx, user.ID)
	if err != nil || streak.Current != 3 {
		t.Errorf("streak = %+v, %v", streak, err)
	}
	progress, err := db.ProgressForCourse(ctx, user.ID, "es-en")
	if err != nil || len(progress) != 1 {
		t.Errorf("progress = %+v, %v", progress, err)
	}
}

func TestSeedWithoutDemoFlagFails(t *testing.T) {
	if _, err := runCmd(t, "seed", "--data", t.TempDir()); err == nil {
		t.Fatal("bare seed succeeded")
	}
}

func TestPackInspectAndValidate(t *testing.T) {
	dir := t.TempDir()

	plain := filepath.Join(dir, "plain.json")
	if err := os.WriteFile(plain, store.DemoPack(), 0o600); err != nil {
		t.Fatal(err)
	}
	enc, _ := zstd.NewWriter(nil)
	packed := filepath.Join(dir, "packed.pack")
	if err := os.WriteFile(packed, enc.EncodeAll(store.DemoPack(), nil), 0o600); err != nil {
		t.Fatal(err)
	}
	enc.Close()

	for _, f := range []string{plain, packed} {
		out, err := runCmd(t, "pack", "inspect", f)
		if err != nil || !strings.Contains(out, "course:  es-en") {
			t.Errorf("inspect %s = %q, %v", f, out, err)
		}
		out, err = runCmd(t, "pack", "validate", f)
		if err != nil || !strings.Contains(out, "header ok") {
			t.Errorf("validate %s = %q, %v", f, out, err)
		}
	}

	bad := filepath.Join(dir, "bad.json")
	os.WriteFile(bad, []byte(`{"format":0}`), 0o600)
	if _, err := runCmd(t, "pack", "validate", bad); err == nil {
		t.Error("validate accepted a bad header")
	}
}

func TestGenRequiresDryRun(t *testing.T) {
	if _, err := runCmd(t, "gen", "--course", "es-en", "--data", t.TempDir()); err == nil {
		t.Fatal("gen without --dry-run succeeded")
	}
}

func TestLoopbackAddr(t *testing.T) {
	cases := []struct {
		in, want string
		bad      bool
	}{
		{in: ":4321", want: "127.0.0.1:4321"},
		{in: "127.0.0.1:4321", want: "127.0.0.1:4321"},
		{in: "localhost:4321", want: "localhost:4321"},
		{in: "[::1]:4321", want: "[::1]:4321"},
		{in: "0.0.0.0:4321", bad: true},
		{in: "192.168.1.4:4321", bad: true},
	}
	for _, tc := range cases {
		got, err := loopbackAddr(tc.in)
		if tc.bad {
			if err == nil {
				t.Errorf("%s: accepted", tc.in)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Errorf("%s = %q, %v, want %q", tc.in, got, err, tc.want)
		}
	}
}

func TestServePrintConfig(t *testing.T) {
	data := t.TempDir()
	t.Setenv("TAMA_LLM_API_KEY", "sk-secret")

	out, err := runCmd(t, "serve", "--print-config", "--data", data, "--addr", ":9999")
	if err != nil {
		t.Fatalf("print-config: %v", err)
	}
	if strings.Contains(out, "sk-secret") {
		t.Error("print-config leaked the API key")
	}
	for _, want := range []string{`addr`, `(flag)`, `"***"`} {
		if !strings.Contains(out, want) {
			t.Errorf("print-config output misses %q:\n%s", want, out)
		}
	}
}
