// Package lang is the compiled-in language registry: every ISO 639-3 code
// plus the constructed and private-use rows the catalog needs (Esperanto,
// Klingon, High Valyrian). The data lives in langdata.go, generated from the
// ISO tables by tools/genlang; nothing here touches the network at runtime.
package lang

import (
	"fmt"
	"sort"
	"strings"
)

//go:generate go run github.com/tamnd/tama/tools/genlang -tsv ../../../tools/genlang/iso-639-3.tsv -o langdata.go

// Text directions. Direction is rtl exactly for the Arab, Hebr, Syrc, and
// Thaa scripts in the seed data, ltr otherwise.
const (
	DirectionLTR = "ltr"
	DirectionRTL = "rtl"
)

// Language is one registry row.
type Language struct {
	Code639_1  string `json:"code639_1,omitempty"` // empty when none exists
	Code639_3  string `json:"code639_3"`           // always present
	Name       string `json:"name"`                // English name
	NativeName string `json:"nativeName"`          // endonym, e.g. 日本語
	Script     string `json:"script,omitempty"`    // ISO 15924, e.g. Latn, Arab
	Direction  string `json:"direction"`           // ltr or rtl
	HasTTS     bool   `json:"hasTTS"`              // the audio pipeline has a voice
}

// Code is the shortest available code, the one course IDs use: 639-1 when it
// exists, 639-3 otherwise.
func (l Language) Code() string {
	if l.Code639_1 != "" {
		return l.Code639_1
	}
	return l.Code639_3
}

// NotFoundError is the typed miss from Lookup.
type NotFoundError struct {
	Code string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("lang: unknown language code %q", e.Code)
}

// index maps every 639-1 code, 639-3 code, and alias (lowercased) onto its
// position in languages. byName is languages sorted by English name.
var (
	index  map[string]int
	byName []Language
)

func init() {
	index = make(map[string]int, 2*len(languages))
	for i := range languages {
		l := &languages[i]
		if l.NativeName == "" {
			l.NativeName = l.Name
		}
		if l.Direction == "" {
			l.Direction = DirectionLTR
		}
		index[l.Code639_3] = i
		if l.Code639_1 != "" {
			index[l.Code639_1] = i
		}
	}
	for alias, code := range aliases {
		index[alias] = index[code]
	}
	byName = make([]Language, len(languages))
	copy(byName, languages)
	sort.Slice(byName, func(i, j int) bool {
		if byName[i].Name != byName[j].Name {
			return byName[i].Name < byName[j].Name
		}
		return byName[i].Code639_3 < byName[j].Code639_3
	})
}

// Lookup resolves a 639-1 code, a 639-3 code, or an alias, case-insensitively
// (zh, zho, and cmn all reach Mandarin Chinese). The error is a *NotFoundError.
func Lookup(code string) (Language, error) {
	i, ok := index[strings.ToLower(code)]
	if !ok {
		return Language{}, &NotFoundError{Code: code}
	}
	return languages[i], nil
}

// All returns every registry row sorted by English name. The slice is a copy;
// callers may reorder it.
func All() []Language {
	out := make([]Language, len(byName))
	copy(out, byName)
	return out
}

// Count is the number of registry rows.
func Count() int {
	return len(languages)
}
