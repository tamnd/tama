package course

import (
	"testing"

	"github.com/tamnd/tama/pkg/course/lang"
)

// TestSeedCatalog is the seed catalog contract: every ID parses, no
// duplicates, every referenced language is in the registry, and all
// from-English targets are present.
func TestSeedCatalog(t *testing.T) {
	all := SeedCourses()
	seen := map[string]bool{}
	for _, c := range all {
		if seen[c.ID] {
			t.Errorf("duplicate seed course %s", c.ID)
		}
		seen[c.ID] = true

		target, base, err := ParseID(c.ID)
		if err != nil {
			t.Errorf("seed course %s does not parse: %v", c.ID, err)
			continue
		}
		if target.Code() != c.TargetLang || base.Code() != c.BaseLang {
			t.Errorf("seed course %s carries langs %s/%s, ID parses to %s/%s",
				c.ID, c.TargetLang, c.BaseLang, target.Code(), base.Code())
		}
		if _, err := lang.Lookup(c.TargetLang); err != nil {
			t.Errorf("seed course %s: target not in registry: %v", c.ID, err)
		}
		if _, err := lang.Lookup(c.BaseLang); err != nil {
			t.Errorf("seed course %s: base not in registry: %v", c.ID, err)
		}
		if c.Order < 1 || c.Title == "" || c.Flag == "" || c.Learners == "" {
			t.Errorf("seed course %s missing picker fields: %+v", c.ID, c)
		}
	}

	wantFromEnglish := []string{
		"es", "fr", "ja", "de", "ko", "it", "zh", "ru", "hi", "pt",
		"ar", "nl", "sv", "el", "ga", "pl", "nb", "tr", "vi", "uk",
		"he", "fi", "da", "cs", "ro", "hu", "cy", "sw", "id", "haw",
		"nv", "gd", "yi", "ht", "zu", "eo", "la", "tlh", "qhv", "ca",
	}
	for _, code := range wantFromEnglish {
		if !seen[code+"-from-en"] {
			t.Errorf("from-English seed target %s missing", code)
		}
	}
	if n := len(SeedCatalog("en")); n != len(wantFromEnglish) {
		t.Errorf("SeedCatalog(en) has %d courses, want %d", n, len(wantFromEnglish))
	}
}

// TestSeedCatalogOrder pins the top of the popularity ladder and the picker
// details called out in the milestone.
func TestSeedCatalogOrder(t *testing.T) {
	en := SeedCatalog("en")
	if len(en) < 2 || en[0].ID != "es-from-en" || en[1].ID != "fr-from-en" {
		t.Fatalf("SeedCatalog(en) does not start with Spanish, French: %+v", en[:2])
	}
	for i, c := range en {
		if c.Order != i+1 {
			t.Errorf("SeedCatalog(en)[%d].Order = %d", i, c.Order)
		}
	}

	byID := map[string]SeedCourse{}
	for _, c := range SeedCourses() {
		byID[c.ID] = c
	}
	if c := byID["nb-from-en"]; c.Title != "Norwegian (Bokmål)" {
		t.Errorf("nb-from-en title = %q, want Norwegian (Bokmål)", c.Title)
	}
	if c := byID["pt-from-en"]; c.Description == "" {
		t.Error("pt-from-en has no Brazilian Portuguese description")
	}
	zh := byID["zh-from-en"]
	if len(zh.ScriptOptions) != 2 || zh.ScriptOptions[0] != "Hans" {
		t.Errorf("zh-from-en script options = %v, want Hans default with Hant", zh.ScriptOptions)
	}
}

// TestSeedCatalogBases checks every non-English base group exists with
// English as a target, and unknown bases return nil.
func TestSeedCatalogBases(t *testing.T) {
	bases := []string{
		"es", "pt", "zh", "ja", "ko", "fr", "de", "ru", "hi", "ar", "vi",
		"tr", "it", "id", "th", "uk", "pl", "cs", "hu", "el", "bn", "te",
		"ta", "nl",
	}
	for _, base := range bases {
		group := SeedCatalog(base)
		if len(group) == 0 {
			t.Errorf("SeedCatalog(%q) is empty", base)
			continue
		}
		if group[0].TargetLang != "en" {
			t.Errorf("SeedCatalog(%q) does not offer English first: %+v", base, group[0])
		}
	}
	if got := SeedCatalog("xx-bogus"); got != nil {
		t.Errorf("SeedCatalog(bogus) = %v, want nil", got)
	}
	if got := SeedCatalog("haw"); got != nil {
		t.Errorf("SeedCatalog(haw) = %v, want nil (no seed courses from Hawaiian)", got)
	}
}
