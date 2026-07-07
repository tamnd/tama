package course

import (
	"errors"
	"testing"

	"github.com/tamnd/tama/pkg/course/lang"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		id           string
		target, base string // want 639-3 codes
	}{
		{"ja-from-en", "jpn", "eng"},
		{"en-from-zh", "eng", "cmn"},
		{"qhv-from-en", "qhv", "eng"},
		{"tlh-from-en", "tlh", "eng"},
		{"haw-from-en", "haw", "eng"},
		{"nb-from-en", "nob", "eng"},
		// Long codes and aliases resolve too; ID gives back the short form.
		{"jpn-from-eng", "jpn", "eng"},
		{"zh-from-cmn", "", ""}, // both sides are Mandarin
		{"en-from-en", "", ""},
		{"en-from-xx", "", ""},
		{"xxq-from-en", "", ""},
		{"es", "", ""},
		{"-from-en", "", ""},
		{"es-from-", "", ""},
	}
	for _, tt := range tests {
		target, base, err := ParseID(tt.id)
		if tt.target == "" {
			if err == nil {
				t.Errorf("ParseID(%q) = %s, %s, nil; want error", tt.id, target.Code639_3, base.Code639_3)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseID(%q): %v", tt.id, err)
			continue
		}
		if target.Code639_3 != tt.target || base.Code639_3 != tt.base {
			t.Errorf("ParseID(%q) = %s, %s; want %s, %s", tt.id, target.Code639_3, base.Code639_3, tt.target, tt.base)
		}
	}
}

// TestParseIDTypedErrors pins the error types callers switch on.
func TestParseIDTypedErrors(t *testing.T) {
	_, _, err := ParseID("en-from-en")
	var same *SamePairError
	if !errors.As(err, &same) {
		t.Errorf("ParseID(en-from-en) err = %T, want *SamePairError", err)
	}

	_, _, err = ParseID("zh-from-cmn")
	if !errors.As(err, &same) {
		t.Errorf("ParseID(zh-from-cmn) err = %T, want *SamePairError (aliases fold)", err)
	}

	_, _, err = ParseID("garbage")
	var invalid *InvalidIDError
	if !errors.As(err, &invalid) {
		t.Errorf("ParseID(garbage) err = %T, want *InvalidIDError", err)
	}

	_, _, err = ParseID("xxq-from-en")
	var nf *lang.NotFoundError
	if !errors.As(err, &nf) {
		t.Errorf("ParseID(xxq-from-en) err = %T, want *lang.NotFoundError", err)
	}
}

// TestIDRoundTrip builds canonical IDs from parsed languages.
func TestIDRoundTrip(t *testing.T) {
	for _, id := range []string{"ja-from-en", "en-from-zh", "qhv-from-en", "haw-from-tlh"} {
		target, base, err := ParseID(id)
		if err != nil {
			t.Fatalf("ParseID(%q): %v", id, err)
		}
		if got := ID(target, base); got != id {
			t.Errorf("ID(ParseID(%q)) = %q", id, got)
		}
	}
	// Non-canonical input normalizes to the short form.
	target, base, err := ParseID("jpn-from-eng")
	if err != nil {
		t.Fatal(err)
	}
	if got := ID(target, base); got != "ja-from-en" {
		t.Errorf("ID(ParseID(jpn-from-eng)) = %q, want ja-from-en", got)
	}
}
