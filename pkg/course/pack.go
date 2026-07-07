package course

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// PackFormatVersion is the pack.json format this build reads and writes.
const PackFormatVersion = 1

// Pack layout on disk: pack.json is the manifest, units/NNN.json holds one
// unit each (levels, lessons, exercises, guidebook, stories), and audio
// blobs are content-addressed under audio/sha256/<aa>/<hash> with their
// metadata in audio/manifest.json.
const (
	manifestFile      = "pack.json"
	unitsDir          = "units"
	audioManifestFile = "audio/manifest.json"
	audioBlobDir      = "audio/sha256"
)

// Manifest is pack.json, the frozen header of a generated pack.
type Manifest struct {
	FormatVersion int    `json:"formatVersion"` // starts at 1
	CourseID      string `json:"courseId"`
	Version       string `json:"version"` // semver
	// GeneratorVersion is the tama version plus prompt-set hash that
	// produced the pack.
	GeneratorVersion string `json:"generatorVersion"`
	// ModelID is the model name reported by the generation endpoint.
	ModelID      string        `json:"modelId"`
	CreatedAt    string        `json:"createdAt"` // RFC 3339
	SectionIndex []SectionMeta `json:"sectionIndex"`
	// ContentHash is sha256 over the sorted (path, file-sha256) list of
	// every file except pack.json itself, so byte-identical trees hash
	// identically regardless of build time.
	ContentHash string `json:"contentHash"`
}

// SectionMeta is one ordered entry of the manifest's section index.
type SectionMeta struct {
	Index int    `json:"index"`
	Title string `json:"title"` // ladder name, "Rookie"
	CEFR  string `json:"cefr"`
	Units int    `json:"units"`
}

// Validate checks the manifest's own fields; cross-file checks live in
// Validate on the whole pack.
func (m Manifest) Validate() error {
	if m.FormatVersion != PackFormatVersion {
		return fmt.Errorf("pack: format version %d, this build reads %d", m.FormatVersion, PackFormatVersion)
	}
	if _, _, err := ParseID(m.CourseID); err != nil {
		return fmt.Errorf("pack: course id: %w", err)
	}
	if _, err := ParseVersion(m.Version); err != nil {
		return err
	}
	if _, err := time.Parse(time.RFC3339, m.CreatedAt); err != nil {
		return fmt.Errorf("pack: createdAt %q is not RFC 3339", m.CreatedAt)
	}
	if len(m.SectionIndex) == 0 {
		return fmt.Errorf("pack: empty section index")
	}
	for i, s := range m.SectionIndex {
		if s.Index != i+1 {
			return fmt.Errorf("pack: section index %d out of order, want %d", s.Index, i+1)
		}
		if !cefrLabels[s.CEFR] {
			return fmt.Errorf("pack: section %d: bad cefr label %q", s.Index, s.CEFR)
		}
		if s.Units < 1 {
			return fmt.Errorf("pack: section %d has no units", s.Index)
		}
	}
	if !strings.HasPrefix(m.ContentHash, "sha256:") {
		return fmt.Errorf("pack: content hash %q is not a sha256 ref", m.ContentHash)
	}
	return nil
}

// UnitCount sums the section index.
func (m Manifest) UnitCount() int {
	n := 0
	for _, s := range m.SectionIndex {
		n += s.Units
	}
	return n
}

// UnitFile is units/NNN.json: everything one unit ships.
type UnitFile struct {
	Unit      Unit       `json:"unit"`
	Levels    []Level    `json:"levels"`
	Lessons   []Lesson   `json:"lessons"`
	Exercises []Exercise `json:"exercises"`
	Guidebook *Guidebook `json:"guidebook,omitempty"`
	Stories   []Story    `json:"stories,omitempty"`
}

// Exercise is the opaque exercise envelope; M4 defines the payload shapes,
// the pack only needs identity and referential integrity.
type Exercise struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Exercise count bands per lesson: outside warnBounds the validator warns,
// outside hardBounds it errors.
var (
	exerciseWarnBounds = [2]int{12, 18}
	exerciseHardBounds = [2]int{8, 25}
)

// AudioManifest is audio/manifest.json.
type AudioManifest struct {
	Entries []AudioEntry `json:"entries"`
}

// AudioEntry describes one content-addressed audio blob.
type AudioEntry struct {
	Hash       string `json:"hash"` // "sha256:<hex>"
	DurationMs int    `json:"durationMs"`
	Voice      string `json:"voice"`
	Kind       string `json:"kind"` // normal, slow, character
	TextHash   string `json:"textHash"`
}

// audioKinds is the closed set for AudioEntry.Kind.
var audioKinds = map[string]bool{"normal": true, "slow": true, "character": true}

// HashHex strips the sha256: prefix, returning an error for anything that
// is not a well-formed ref.
func HashHex(ref string) (string, error) {
	hexPart, ok := strings.CutPrefix(ref, "sha256:")
	if !ok || len(hexPart) != 64 {
		return "", fmt.Errorf("pack: %q is not a sha256 ref", ref)
	}
	if _, err := hex.DecodeString(hexPart); err != nil {
		return "", fmt.Errorf("pack: %q is not a sha256 ref", ref)
	}
	return hexPart, nil
}

// BlobPath is the pack-relative path of an audio ref's blob.
func BlobPath(ref string) (string, error) {
	hexPart, err := HashHex(ref)
	if err != nil {
		return "", err
	}
	return audioBlobDir + "/" + hexPart[:2] + "/" + hexPart, nil
}

// Version is a parsed semver. PATCH is for audio and typo fixes with no ID
// changes, MINOR adds units, levels, or stories, MAJOR removes or reshapes
// existing IDs.
type Version struct {
	Major, Minor, Patch int
}

// ParseVersion reads strict MAJOR.MINOR.PATCH.
func ParseVersion(s string) (Version, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("pack: version %q is not MAJOR.MINOR.PATCH", s)
	}
	var nums [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 || (len(p) > 1 && p[0] == '0') {
			return Version{}, fmt.Errorf("pack: version %q is not MAJOR.MINOR.PATCH", s)
		}
		nums[i] = n
	}
	return Version{nums[0], nums[1], nums[2]}, nil
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Compare orders versions: -1, 0, or 1.
func (v Version) Compare(o Version) int {
	for _, d := range [3]int{v.Major - o.Major, v.Minor - o.Minor, v.Patch - o.Patch} {
		if d < 0 {
			return -1
		}
		if d > 0 {
			return 1
		}
	}
	return 0
}

// DowngradeError is the typed refusal of a MAJOR downgrade.
type DowngradeError struct {
	Installed, Incoming Version
}

func (e *DowngradeError) Error() string {
	return fmt.Sprintf("pack: refusing major downgrade from %s to %s", e.Installed, e.Incoming)
}

// CheckUpgrade decides whether incoming may replace installed. Anything
// except a MAJOR downgrade goes through; reinstalling the same or an older
// minor of the same major is allowed for rollbacks.
func CheckUpgrade(installed, incoming string) error {
	have, err := ParseVersion(installed)
	if err != nil {
		return err
	}
	want, err := ParseVersion(incoming)
	if err != nil {
		return err
	}
	if want.Major < have.Major {
		return &DowngradeError{Installed: have, Incoming: want}
	}
	return nil
}

// Pack is one loaded pack directory.
type Pack struct {
	Dir      string
	Manifest Manifest
	Units    []UnitFile
	Audio    AudioManifest
}

// LoadManifest reads and parses pack.json without validating it.
func LoadManifest(dir string) (Manifest, error) {
	var m Manifest
	raw, err := os.ReadFile(filepath.Join(dir, manifestFile))
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		return m, fmt.Errorf("pack: %s: %w", manifestFile, err)
	}
	return m, nil
}

// LoadPack reads a pack directory into memory: manifest, every unit file in
// name order, and the audio manifest. It parses but does not validate; run
// Validate for the full contract.
func LoadPack(dir string) (*Pack, error) {
	m, err := LoadManifest(dir)
	if err != nil {
		return nil, err
	}
	p := &Pack{Dir: dir, Manifest: m}

	names, err := unitFileNames(dir)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		raw, err := os.ReadFile(filepath.Join(dir, unitsDir, name))
		if err != nil {
			return nil, err
		}
		var u UnitFile
		if err := json.Unmarshal(raw, &u); err != nil {
			return nil, fmt.Errorf("pack: %s/%s: %w", unitsDir, name, err)
		}
		p.Units = append(p.Units, u)
	}

	raw, err := os.ReadFile(filepath.Join(dir, audioManifestFile))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &p.Audio); err != nil {
		return nil, fmt.Errorf("pack: %s: %w", audioManifestFile, err)
	}
	return p, nil
}

// unitFileNames lists units/*.json sorted by name.
func unitFileNames(dir string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(dir, unitsDir))
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// ComputeContentHash hashes the pack tree: sha256 over the sorted
// (path, file-sha256) list of every file under dir except pack.json itself.
// Timestamps and permissions never feed the hash, so byte-identical trees
// match regardless of when or where they were built.
func ComputeContentHash(dir string) (string, error) {
	type fileHash struct {
		path, sum string
	}
	var files []fileHash
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == manifestFile {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		files = append(files, fileHash{rel, hex.EncodeToString(h.Sum(nil))})
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })

	h := sha256.New()
	for _, f := range files {
		fmt.Fprintf(h, "%s\t%s\n", f.path, f.sum)
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}
