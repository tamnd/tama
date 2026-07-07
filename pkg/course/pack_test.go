package course

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixturePack = "testdata/pack-es-from-en-mini"

// TestFixturePackShape pins the mini pack's exact counts: 2 sections,
// 3 units, 8 levels, 1 story, 1 guidebook per unit.
func TestFixturePackShape(t *testing.T) {
	p, err := LoadPack(fixturePack)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(p.Manifest.SectionIndex); got != 2 {
		t.Errorf("sections = %d, want 2", got)
	}
	if got := len(p.Units); got != 3 {
		t.Fatalf("units = %d, want 3", got)
	}

	levels, stories, guidebooks := 0, 0, 0
	kinds := map[NodeKind]int{}
	for _, u := range p.Units {
		levels += len(u.Levels)
		stories += len(u.Stories)
		if u.Guidebook != nil {
			guidebooks++
		}
		for _, l := range u.Levels {
			kinds[l.Kind]++
		}
	}
	if levels != 8 {
		t.Errorf("levels = %d, want 8", levels)
	}
	if stories != 1 {
		t.Errorf("stories = %d, want 1", stories)
	}
	if guidebooks != 3 {
		t.Errorf("guidebooks = %d, want one per unit", guidebooks)
	}
	for _, k := range []NodeKind{KindLesson, KindChest, KindStory, KindHard, KindReview, KindCheckpoint} {
		if kinds[k] == 0 {
			t.Errorf("fixture has no %s node", k)
		}
	}

	st := p.Units[0].Stories[0]
	if len(st.Lines) != 12 || len(st.Characters) != 3 || len(st.Exercises) != 3 {
		t.Errorf("story has %d lines, %d characters, %d exercises; want 12, 3, 3",
			len(st.Lines), len(st.Characters), len(st.Exercises))
	}
	if st.XPReward != XPStory {
		t.Errorf("story xp = %d, want %d", st.XPReward, XPStory)
	}
}

// TestFixturePackValidates wants the checked-in pack fully green: no errors
// and no warnings.
func TestFixturePackValidates(t *testing.T) {
	r := Validate(fixturePack)
	for _, e := range r.Errors {
		t.Errorf("error: %s", e)
	}
	for _, w := range r.Warnings {
		t.Errorf("warning: %s", w)
	}
}

// copyFixture clones the fixture pack into a temp dir the test may corrupt.
func copyFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS(fixturePack)); err != nil {
		t.Fatal(err)
	}
	return dir
}

// editJSON loads a JSON file into a generic map, hands it to mut, and
// writes it back.
func editJSON(t *testing.T, path string, mut func(map[string]any)) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	mut(doc)
	out, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestValidateCatchesCorruption seeds one defect per case into a copy of
// the fixture and wants a matching error.
func TestValidateCatchesCorruption(t *testing.T) {
	blobRel := filepath.FromSlash("audio/sha256/47/476e8b59e0e8ac3cd25a9079614188944feffc85a4942a1fd5009e6f052966d9")
	tests := []struct {
		name    string
		corrupt func(t *testing.T, dir string)
		want    string // substring of an expected error
	}{
		{
			"bad blob hash",
			func(t *testing.T, dir string) {
				path := filepath.Join(dir, blobRel)
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				raw[len(raw)-1] ^= 0xff
				if err := os.WriteFile(path, raw, 0o644); err != nil {
					t.Fatal(err)
				}
			},
			"blob hashes to",
		},
		{
			"missing audio blob",
			func(t *testing.T, dir string) {
				if err := os.Remove(filepath.Join(dir, blobRel)); err != nil {
					t.Fatal(err)
				}
			},
			"has no blob",
		},
		{
			"dangling exercise ref",
			func(t *testing.T, dir string) {
				editJSON(t, filepath.Join(dir, "units", "002.json"), func(doc map[string]any) {
					lesson := doc["lessons"].([]any)[0].(map[string]any)
					lesson["exerciseIds"].([]any)[0] = "ex_zz_99"
				})
			},
			"unknown exercise ex_zz_99",
		},
		{
			"malformed guidebook markdown",
			func(t *testing.T, dir string) {
				editJSON(t, filepath.Join(dir, "units", "003.json"), func(doc map[string]any) {
					gb := doc["guidebook"].(map[string]any)
					tip := gb["tipSections"].([]any)[0].(map[string]any)
					tip["bodyMarkdown"] = "# Mi and mis\nSee [here](https://example.com)."
				})
			},
			"not allowed in guidebook markdown",
		},
		{
			"dangling story node",
			func(t *testing.T, dir string) {
				editJSON(t, filepath.Join(dir, "units", "001.json"), func(doc map[string]any) {
					level := doc["levels"].([]any)[2].(map[string]any)
					level["payload"].(map[string]any)["storyId"] = "st_missing"
				})
			},
			`story "st_missing"`,
		},
		{
			"content edit without rehash",
			func(t *testing.T, dir string) {
				editJSON(t, filepath.Join(dir, "units", "001.json"), func(doc map[string]any) {
					ex := doc["exercises"].([]any)[0].(map[string]any)
					ex["payload"].(map[string]any)["answer"] = "el cafe"
				})
			},
			"content hash",
		},
		{
			"section promises more units",
			func(t *testing.T, dir string) {
				editJSON(t, filepath.Join(dir, "pack.json"), func(doc map[string]any) {
					sec := doc["sectionIndex"].([]any)[1].(map[string]any)
					sec["units"] = 2.0
				})
			},
			"promises",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := copyFixture(t)
			tt.corrupt(t, dir)
			r := Validate(dir)
			if r.OK() {
				t.Fatal("corrupted pack validated clean")
			}
			found := false
			for _, e := range r.Errors {
				if strings.Contains(e, tt.want) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("no error containing %q, got %q", tt.want, r.Errors)
			}
		})
	}
}

// TestContainerRoundTrip encodes the fixture, decodes it elsewhere, and
// wants an identical content hash and a clean validation; encoding twice
// must give identical bytes.
func TestContainerRoundTrip(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	if err := WriteContainer(fixturePack, &buf1); err != nil {
		t.Fatal(err)
	}
	if err := WriteContainer(fixturePack, &buf2); err != nil {
		t.Fatal(err)
	}
	if sha256.Sum256(buf1.Bytes()) != sha256.Sum256(buf2.Bytes()) {
		t.Error("container encoding is not deterministic")
	}

	dir := t.TempDir()
	if err := ReadContainer(bytes.NewReader(buf1.Bytes()), dir); err != nil {
		t.Fatal(err)
	}
	want, err := ComputeContentHash(fixturePack)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ComputeContentHash(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("content hash after round trip = %s, want %s", got, want)
	}
	if r := Validate(dir); !r.OK() {
		t.Errorf("round-tripped pack fails validation: %q", r.Errors)
	}

	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if m.ContentHash != want {
		t.Errorf("manifest hash %s, tree hash %s", m.ContentHash, want)
	}
}

// TestReadContainerRejectsEscapes blocks path traversal out of the target.
func TestReadContainerRejectsEscapes(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteContainer(fixturePack, &buf); err != nil {
		t.Fatal(err)
	}
	// Splice a hostile name into a fresh archive by re-tarring one entry.
	evil := t.TempDir()
	if err := os.WriteFile(filepath.Join(evil, "ok"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	var hostile bytes.Buffer
	if err := WriteContainer(evil, &hostile); err != nil {
		t.Fatal(err)
	}
	raw := bytes.Replace(hostile.Bytes(), []byte("ok"), []byte(".."), 1)
	if err := ReadContainer(bytes.NewReader(raw), t.TempDir()); err == nil {
		t.Error("hostile entry extracted without error")
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		in   string
		want Version
		ok   bool
	}{
		{"1.0.0", Version{1, 0, 0}, true},
		{"1.2.10", Version{1, 2, 10}, true},
		{"0.1.0", Version{0, 1, 0}, true},
		{"1.2", Version{}, false},
		{"1.2.3.4", Version{}, false},
		{"v1.2.3", Version{}, false},
		{"1.02.3", Version{}, false},
		{"1.2.x", Version{}, false},
		{"", Version{}, false},
	}
	for _, tt := range tests {
		got, err := ParseVersion(tt.in)
		if (err == nil) != tt.ok || got != tt.want {
			t.Errorf("ParseVersion(%q) = %v, %v; want %v, ok=%v", tt.in, got, err, tt.want, tt.ok)
		}
	}
	if (Version{2, 0, 0}).Compare(Version{1, 9, 9}) != 1 ||
		(Version{1, 2, 3}).Compare(Version{1, 2, 3}) != 0 ||
		(Version{1, 2, 3}).Compare(Version{1, 3, 0}) != -1 {
		t.Error("Version.Compare misorders")
	}
}

// TestCheckUpgrade allows everything but a MAJOR downgrade.
func TestCheckUpgrade(t *testing.T) {
	for _, ok := range []struct{ installed, incoming string }{
		{"1.0.0", "1.0.1"}, {"1.0.0", "1.1.0"}, {"1.0.0", "2.0.0"},
		{"1.2.0", "1.1.0"}, {"1.0.0", "1.0.0"},
	} {
		if err := CheckUpgrade(ok.installed, ok.incoming); err != nil {
			t.Errorf("CheckUpgrade(%s, %s) = %v, want nil", ok.installed, ok.incoming, err)
		}
	}
	err := CheckUpgrade("2.0.0", "1.9.9")
	var dg *DowngradeError
	if !errors.As(err, &dg) {
		t.Errorf("CheckUpgrade major downgrade err = %T, want *DowngradeError", err)
	}
}

func TestManifestValidate(t *testing.T) {
	good, err := LoadManifest(fixturePack)
	if err != nil {
		t.Fatal(err)
	}
	if err := good.Validate(); err != nil {
		t.Fatalf("fixture manifest: %v", err)
	}
	tests := []struct {
		name string
		mut  func(*Manifest)
	}{
		{"wrong format", func(m *Manifest) { m.FormatVersion = 2 }},
		{"bad course id", func(m *Manifest) { m.CourseID = "es-from-es" }},
		{"bad version", func(m *Manifest) { m.Version = "1.0" }},
		{"bad time", func(m *Manifest) { m.CreatedAt = "yesterday" }},
		{"no sections", func(m *Manifest) { m.SectionIndex = nil }},
		{"gap in sections", func(m *Manifest) { m.SectionIndex[1].Index = 3 }},
		{"bad hash", func(m *Manifest) { m.ContentHash = "md5:abc" }},
	}
	for _, tt := range tests {
		m := good
		m.SectionIndex = append([]SectionMeta(nil), good.SectionIndex...)
		tt.mut(&m)
		if err := m.Validate(); err == nil {
			t.Errorf("%s: want error", tt.name)
		}
	}
}

func TestBlobPath(t *testing.T) {
	ref := "sha256:476e8b59e0e8ac3cd25a9079614188944feffc85a4942a1fd5009e6f052966d9"
	got, err := BlobPath(ref)
	if err != nil {
		t.Fatal(err)
	}
	want := "audio/sha256/47/476e8b59e0e8ac3cd25a9079614188944feffc85a4942a1fd5009e6f052966d9"
	if got != want {
		t.Errorf("BlobPath = %s, want %s", got, want)
	}
	for _, bad := range []string{"476e8b59", "sha256:zz", "sha256:abc", "md5:476e8b59e0e8ac3cd25a9079614188944feffc85a4942a1fd5009e6f052966d9"} {
		if _, err := BlobPath(bad); err == nil {
			t.Errorf("BlobPath(%q): want error", bad)
		}
	}
}
