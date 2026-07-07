package course

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Report is the outcome of validating one pack: hard errors that make the
// pack unusable and warnings the generator should look at.
type Report struct {
	Errors   []string
	Warnings []string
}

// OK reports whether the pack passed with no errors; warnings do not fail
// a pack.
func (r *Report) OK() bool { return len(r.Errors) == 0 }

func (r *Report) errf(format string, args ...any) {
	r.Errors = append(r.Errors, fmt.Sprintf(format, args...))
}

func (r *Report) warnf(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

// Validate checks a pack directory against the full contract: manifest
// fields, unit file numbering, structural rules on every row, referential
// integrity (lessons to exercises, nodes to stories, audio refs to manifest
// entries to blobs), guidebook markdown restrictions, blob hashes, and the
// manifest content hash over the tree.
func Validate(dir string) *Report {
	r := &Report{}

	m, err := LoadManifest(dir)
	if err != nil {
		r.errf("%v", err)
		return r
	}
	if err := m.Validate(); err != nil {
		r.errf("%s: %v", manifestFile, err)
	}

	names, err := unitFileNames(dir)
	if err != nil {
		r.errf("%v", err)
		return r
	}
	if want := m.UnitCount(); len(names) != want {
		r.errf("%s: section index promises %d units, %s/ has %d files", manifestFile, want, unitsDir, len(names))
	}
	for i, name := range names {
		if want := fmt.Sprintf("%03d.json", i+1); name != want {
			r.errf("%s/%s: unit files must be numbered sequentially, want %s", unitsDir, name, want)
		}
	}

	p, err := LoadPack(dir)
	if err != nil {
		r.errf("%v", err)
		return r
	}

	ids := map[string]string{} // id -> kind of thing, for cross-file uniqueness
	claim := func(id, what, where string) {
		if prev, dup := ids[id]; dup {
			r.errf("%s: %s id %q already used by a %s", where, what, id, prev)
			return
		}
		ids[id] = what
	}

	audioRefs := map[string][]string{} // ref -> where it was seen
	stories := map[string]bool{}
	storyRefs := map[string][]string{} // story id -> nodes linking to it

	sectionUnits := map[int]int{}
	for i, u := range p.Units {
		where := fmt.Sprintf("%s/%03d.json", unitsDir, i+1)
		validateUnitFile(r, u, where, claim, audioRefs, stories, storyRefs, i == len(p.Units)-1)
		if u.Unit.SectionIndex >= 1 && u.Unit.SectionIndex <= len(m.SectionIndex) {
			sectionUnits[u.Unit.SectionIndex]++
			if u.Unit.Index != sectionUnits[u.Unit.SectionIndex] {
				r.errf("%s: unit %s has index %d, want %d within section %d",
					where, u.Unit.ID, u.Unit.Index, sectionUnits[u.Unit.SectionIndex], u.Unit.SectionIndex)
			}
		} else if u.Unit.SectionIndex != 0 {
			r.errf("%s: unit %s points at section %d, manifest has %d sections",
				where, u.Unit.ID, u.Unit.SectionIndex, len(m.SectionIndex))
		}
	}
	for _, s := range m.SectionIndex {
		if sectionUnits[s.Index] != s.Units {
			r.errf("%s: section %d promises %d units, unit files deliver %d",
				manifestFile, s.Index, s.Units, sectionUnits[s.Index])
		}
	}
	for id, nodes := range storyRefs {
		if !stories[id] {
			r.errf("%s: story %q referenced by %s is not in the pack", unitsDir, id, strings.Join(nodes, ", "))
		}
	}

	validateAudio(r, dir, p.Audio, audioRefs)

	if r.OK() {
		sum, err := ComputeContentHash(dir)
		if err != nil {
			r.errf("%v", err)
		} else if sum != m.ContentHash {
			r.errf("%s: content hash %s does not match tree %s", manifestFile, m.ContentHash, sum)
		}
	}
	return r
}

// validateUnitFile checks one unit file's rows and referential integrity.
func validateUnitFile(r *Report, u UnitFile, where string, claim func(id, what, where string),
	audioRefs map[string][]string, stories map[string]bool, storyRefs map[string][]string, lastUnit bool) {

	if err := u.Unit.Validate(); err != nil {
		r.errf("%s: %v", where, err)
	}
	claim(u.Unit.ID, "unit", where)

	exercises := map[string]bool{}
	for _, ex := range u.Exercises {
		if ex.ID == "" || ex.Type == "" {
			r.errf("%s: exercise missing id or type", where)
			continue
		}
		claim(ex.ID, "exercise", where)
		exercises[ex.ID] = true
	}

	levels := map[string]*Level{}
	for j := range u.Levels {
		l := &u.Levels[j]
		if err := l.Validate(); err != nil {
			r.errf("%s: %v", where, err)
			continue
		}
		claim(l.ID, "level", where)
		if l.UnitID != u.Unit.ID {
			r.errf("%s: level %s belongs to unit %s, file is unit %s", where, l.ID, l.UnitID, u.Unit.ID)
		}
		if l.Index != j+1 {
			r.errf("%s: level %s has index %d, want %d", where, l.ID, l.Index, j+1)
		}
		levels[l.ID] = l
		validateNodePayload(r, *l, where, storyRefs)
	}
	if len(u.Levels) != u.Unit.LevelCount {
		r.errf("%s: unit %s says %d levels, file has %d", where, u.Unit.ID, u.Unit.LevelCount, len(u.Levels))
	}
	if len(u.Levels) > 0 && u.Levels[len(u.Levels)-1].Kind != KindCheckpoint {
		r.errf("%s: unit %s does not end with a checkpoint", where, u.Unit.ID)
	}

	lessonsPerLevel := map[string]int{}
	for _, ls := range u.Lessons {
		if err := ls.Validate(); err != nil {
			r.errf("%s: %v", where, err)
			continue
		}
		claim(ls.ID, "lesson", where)
		if _, ok := levels[ls.LevelID]; !ok {
			r.errf("%s: lesson %s points at unknown level %s", where, ls.ID, ls.LevelID)
			continue
		}
		lessonsPerLevel[ls.LevelID]++
		if ls.Index != lessonsPerLevel[ls.LevelID] {
			r.errf("%s: lesson %s has index %d, want %d within level %s",
				where, ls.ID, ls.Index, lessonsPerLevel[ls.LevelID], ls.LevelID)
		}
		n := len(ls.ExerciseIDs)
		if n < exerciseHardBounds[0] || n > exerciseHardBounds[1] {
			r.errf("%s: lesson %s has %d exercises, want %d to %d",
				where, ls.ID, n, exerciseHardBounds[0], exerciseHardBounds[1])
		} else if n < exerciseWarnBounds[0] || n > exerciseWarnBounds[1] {
			r.warnf("%s: lesson %s has %d exercises, %d to %d reads better",
				where, ls.ID, n, exerciseWarnBounds[0], exerciseWarnBounds[1])
		}
		for _, exID := range ls.ExerciseIDs {
			if !exercises[exID] {
				r.errf("%s: lesson %s references unknown exercise %s", where, ls.ID, exID)
			}
		}
	}
	for id, l := range levels {
		if lessonsPerLevel[id] != l.LessonCount {
			r.errf("%s: level %s says %d lessons, file has %d", where, id, l.LessonCount, lessonsPerLevel[id])
		}
	}

	switch {
	case u.Guidebook == nil:
		// Only the final review unit may go without a guidebook.
		if !lastUnit {
			r.errf("%s: unit %s has no guidebook", where, u.Unit.ID)
		}
	default:
		gb := u.Guidebook
		claim(gb.ID, "guidebook", where)
		if gb.UnitID != u.Unit.ID {
			r.errf("%s: guidebook %s belongs to unit %s, file is unit %s", where, gb.ID, gb.UnitID, u.Unit.ID)
		}
		if u.Unit.GuidebookID != gb.ID {
			r.errf("%s: unit %s points at guidebook %q, file carries %q", where, u.Unit.ID, u.Unit.GuidebookID, gb.ID)
		}
		if gb.Title == "" || len(gb.KeyPhrases) == 0 {
			r.errf("%s: guidebook %s needs a title and key phrases", where, gb.ID)
		}
		for _, kp := range gb.KeyPhrases {
			if kp.TargetText == "" || kp.BaseText == "" {
				r.errf("%s: guidebook %s has an empty key phrase", where, gb.ID)
			}
			if kp.AudioRef != "" {
				audioRefs[kp.AudioRef] = append(audioRefs[kp.AudioRef], "guidebook "+gb.ID)
			}
		}
		for _, ts := range gb.TipSections {
			if err := validateTipMarkdown(ts.BodyMarkdown); err != nil {
				r.errf("%s: guidebook %s, section %q: %v", where, gb.ID, ts.Heading, err)
			}
		}
	}

	for _, st := range u.Stories {
		if err := st.Validate(); err != nil {
			r.errf("%s: %v", where, err)
			continue
		}
		claim(st.ID, "story", where)
		if st.UnitID != u.Unit.ID {
			r.errf("%s: story %s belongs to unit %s, file is unit %s", where, st.ID, st.UnitID, u.Unit.ID)
		}
		stories[st.ID] = true
		for _, line := range st.Lines {
			if line.AudioRef != "" {
				audioRefs[line.AudioRef] = append(audioRefs[line.AudioRef], "story "+st.ID)
			}
		}
	}
}

// validateNodePayload parses the kind-discriminated payload of one level.
func validateNodePayload(r *Report, l Level, where string, storyRefs map[string][]string) {
	switch l.Kind {
	case KindLesson, KindHard:
		var pl LessonPayload
		if err := json.Unmarshal(l.Payload, &pl); err != nil || pl.XP < 1 {
			r.errf("%s: level %s needs a lesson payload with xp, got %s", where, l.ID, l.Payload)
		}
	case KindChest:
		var pl ChestPayload
		if err := json.Unmarshal(l.Payload, &pl); err != nil || pl.Gems < 1 {
			r.errf("%s: level %s needs a chest payload with gems, got %s", where, l.ID, l.Payload)
		}
	case KindStory:
		var pl StoryPayload
		if err := json.Unmarshal(l.Payload, &pl); err != nil || pl.StoryID == "" {
			r.errf("%s: level %s needs a story payload with storyId, got %s", where, l.ID, l.Payload)
			return
		}
		storyRefs[pl.StoryID] = append(storyRefs[pl.StoryID], l.ID)
	}
}

// validateTipMarkdown enforces the guidebook body restriction: bold, italic,
// and tables only. Headings, links, images, code, and raw HTML are out.
func validateTipMarkdown(body string) error {
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("empty body markdown")
	}
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			return fmt.Errorf("headings are not allowed in guidebook markdown")
		}
	}
	for _, bad := range []struct{ marker, what string }{
		{"](", "links"},
		{"![", "images"},
		{"`", "code"},
		{"<", "raw HTML"},
	} {
		if strings.Contains(body, bad.marker) {
			return fmt.Errorf("%s are not allowed in guidebook markdown", bad.what)
		}
	}
	return nil
}

// validateAudio checks refs against the manifest and the manifest against
// the blobs on disk, rehashing every blob.
func validateAudio(r *Report, dir string, am AudioManifest, refs map[string][]string) {
	entries := map[string]bool{}
	for _, e := range am.Entries {
		hexPart, err := HashHex(e.Hash)
		if err != nil {
			r.errf("%s: %v", audioManifestFile, err)
			continue
		}
		if entries[e.Hash] {
			r.errf("%s: duplicate entry %s", audioManifestFile, e.Hash)
		}
		entries[e.Hash] = true
		if !audioKinds[e.Kind] {
			r.errf("%s: entry %s has kind %q, want normal, slow, or character", audioManifestFile, e.Hash, e.Kind)
		}
		if e.DurationMs < 1 || e.Voice == "" {
			r.errf("%s: entry %s needs a duration and voice", audioManifestFile, e.Hash)
		}

		rel, _ := BlobPath(e.Hash)
		f, err := os.Open(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil {
			r.errf("%s: entry %s has no blob at %s", audioManifestFile, e.Hash, rel)
			continue
		}
		h := sha256.New()
		_, err = io.Copy(h, f)
		f.Close()
		if err != nil {
			r.errf("%s: %v", rel, err)
			continue
		}
		if got := hex.EncodeToString(h.Sum(nil)); got != hexPart {
			r.errf("%s: blob hashes to sha256:%s, manifest says %s", rel, got, e.Hash)
		}
	}
	for ref, seen := range refs {
		if !entries[ref] {
			r.errf("audio ref %s (%s) has no manifest entry", ref, strings.Join(seen, ", "))
		}
	}
}
