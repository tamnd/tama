package lang

import (
	"errors"
	"sort"
	"testing"
)

// TestRegistryCount pins the generated row count: the ISO table minus the
// four special codes and the folded zho row, plus the private-use qhv row.
// Regenerating from a newer table moves this number on purpose.
func TestRegistryCount(t *testing.T) {
	if got := Count(); got != 7919 {
		t.Errorf("Count() = %d, want 7919", got)
	}
}

// TestRegistrySpotChecks pins name, endonym, script, and direction for the
// languages the seed catalog leans on.
func TestRegistrySpotChecks(t *testing.T) {
	tests := []struct {
		code      string
		name      string
		native    string
		script    string
		direction string
	}{
		{"ja", "Japanese", "日本語", "Jpan", DirectionLTR},
		{"zh", "Chinese", "中文", "Hani", DirectionLTR},
		{"ko", "Korean", "한국어", "Hang", DirectionLTR},
		{"ar", "Arabic", "العربية", "Arab", DirectionRTL},
		{"he", "Hebrew", "עברית", "Hebr", DirectionRTL},
		{"hi", "Hindi", "हिन्दी", "Deva", DirectionLTR},
		{"vi", "Vietnamese", "Tiếng Việt", "Latn", DirectionLTR},
		{"cy", "Welsh", "Cymraeg", "Latn", DirectionLTR},
		{"gd", "Scottish Gaelic", "Gàidhlig", "Latn", DirectionLTR},
		{"nv", "Navajo", "Diné bizaad", "Latn", DirectionLTR},
		{"haw", "Hawaiian", "ʻŌlelo Hawaiʻi", "Latn", DirectionLTR},
		{"ht", "Haitian Creole", "Kreyòl ayisyen", "Latn", DirectionLTR},
		{"tlh", "Klingon", "tlhIngan Hol", "Latn", DirectionLTR},
		{"qhv", "High Valyrian", "Valyrio", "Latn", DirectionLTR},
		{"es", "Spanish", "Español", "Latn", DirectionLTR},
		{"ru", "Russian", "Русский", "Cyrl", DirectionLTR},
		{"el", "Greek", "Ελληνικά", "Grek", DirectionLTR},
		{"nb", "Norwegian (Bokmål)", "Norsk (bokmål)", "Latn", DirectionLTR},
		{"yi", "Yiddish", "ייִדיש", "Hebr", DirectionRTL},
		{"th", "Thai", "ไทย", "Thai", DirectionLTR},
	}
	for _, tt := range tests {
		l, err := Lookup(tt.code)
		if err != nil {
			t.Errorf("Lookup(%q): %v", tt.code, err)
			continue
		}
		if l.Name != tt.name || l.NativeName != tt.native || l.Script != tt.script || l.Direction != tt.direction {
			t.Errorf("Lookup(%q) = {%s %s %s %s}, want {%s %s %s %s}",
				tt.code, l.Name, l.NativeName, l.Script, l.Direction,
				tt.name, tt.native, tt.script, tt.direction)
		}
	}
}

// TestDirection asserts the RTL rule: ar, he, yi are rtl; en, ja, ru are ltr.
func TestDirection(t *testing.T) {
	for _, code := range []string{"ar", "he", "yi", "fa", "ur", "dv"} {
		if l, err := Lookup(code); err != nil || l.Direction != DirectionRTL {
			t.Errorf("Lookup(%q).Direction = %q (%v), want rtl", code, l.Direction, err)
		}
	}
	for _, code := range []string{"en", "ja", "ru"} {
		if l, err := Lookup(code); err != nil || l.Direction != DirectionLTR {
			t.Errorf("Lookup(%q).Direction = %q (%v), want ltr", code, l.Direction, err)
		}
	}
}

// TestLookupAliases checks 639-1, 639-3, aliases, and case folding all reach
// the same row, and that Mandarin owns zh, zho, and cmn.
func TestLookupAliases(t *testing.T) {
	for _, code := range []string{"zh", "zho", "cmn", "ZH", "Cmn"} {
		l, err := Lookup(code)
		if err != nil {
			t.Fatalf("Lookup(%q): %v", code, err)
		}
		if l.Code639_3 != "cmn" || l.Name != "Chinese" {
			t.Errorf("Lookup(%q) = %s (%s), want Chinese (cmn)", code, l.Name, l.Code639_3)
		}
	}
	ja1, _ := Lookup("ja")
	ja3, _ := Lookup("jpn")
	if ja1 != ja3 {
		t.Errorf("ja and jpn resolve to different rows: %+v vs %+v", ja1, ja3)
	}
}

// TestLookupNotFound wants a typed error for garbage input.
func TestLookupNotFound(t *testing.T) {
	_, err := Lookup("xx-bogus")
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("Lookup(garbage) err = %v, want *NotFoundError", err)
	}
	if nf.Code != "xx-bogus" {
		t.Errorf("NotFoundError.Code = %q", nf.Code)
	}
}

// TestCode picks the shortest available code.
func TestCode(t *testing.T) {
	ja, _ := Lookup("jpn")
	if ja.Code() != "ja" {
		t.Errorf("jpn Code() = %q, want ja", ja.Code())
	}
	haw, _ := Lookup("haw")
	if haw.Code() != "haw" {
		t.Errorf("haw Code() = %q, want haw", haw.Code())
	}
}

// TestAllSorted wants All() sorted by English name with defaults filled.
func TestAllSorted(t *testing.T) {
	all := All()
	if len(all) != Count() {
		t.Fatalf("All() = %d rows, want %d", len(all), Count())
	}
	if !sort.SliceIsSorted(all, func(i, j int) bool {
		if all[i].Name != all[j].Name {
			return all[i].Name < all[j].Name
		}
		return all[i].Code639_3 < all[j].Code639_3
	}) {
		t.Error("All() is not sorted by English name")
	}
	for _, l := range all {
		if l.NativeName == "" || l.Direction == "" || l.Code639_3 == "" {
			t.Fatalf("row %q missing defaults: %+v", l.Code639_3, l)
		}
	}
}
