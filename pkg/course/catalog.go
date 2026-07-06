// Package course defines the course model: languages, the catalog of
// base-to-target pairs, and (in later milestones) the path structure loaded
// from course packs. Any base to any target is a valid course; the seed list
// below is just what the catalog screen surfaces first.
package course

import "fmt"

// Language is one entry in the registry.
type Language struct {
	Code   string `json:"code"`   // ISO 639-1 where it exists
	Name   string `json:"name"`   // English name
	Native string `json:"native"` // endonym, shown on flags/cards
	RTL    bool   `json:"rtl,omitempty"`
}

// Course is a base-to-target pair. ID is "<target>-from-<base>".
type Course struct {
	ID     string   `json:"id"`
	Base   Language `json:"base"`
	Target Language `json:"target"`
}

// Languages is the seed registry. The full ISO set lands with M3; these cover
// the seed catalog.
var Languages = []Language{
	{Code: "en", Name: "English", Native: "English"},
	{Code: "es", Name: "Spanish", Native: "Español"},
	{Code: "fr", Name: "French", Native: "Français"},
	{Code: "ja", Name: "Japanese", Native: "日本語"},
	{Code: "de", Name: "German", Native: "Deutsch"},
	{Code: "ko", Name: "Korean", Native: "한국어"},
	{Code: "it", Name: "Italian", Native: "Italiano"},
	{Code: "zh", Name: "Chinese", Native: "中文"},
	{Code: "ru", Name: "Russian", Native: "Русский"},
	{Code: "hi", Name: "Hindi", Native: "हिन्दी"},
	{Code: "pt", Name: "Portuguese", Native: "Português"},
	{Code: "ar", Name: "Arabic", Native: "العربية", RTL: true},
	{Code: "nl", Name: "Dutch", Native: "Nederlands"},
	{Code: "sv", Name: "Swedish", Native: "Svenska"},
	{Code: "el", Name: "Greek", Native: "Ελληνικά"},
	{Code: "ga", Name: "Irish", Native: "Gaeilge"},
	{Code: "pl", Name: "Polish", Native: "Polski"},
	{Code: "no", Name: "Norwegian", Native: "Norsk"},
	{Code: "tr", Name: "Turkish", Native: "Türkçe"},
	{Code: "vi", Name: "Vietnamese", Native: "Tiếng Việt"},
	{Code: "uk", Name: "Ukrainian", Native: "Українська"},
	{Code: "he", Name: "Hebrew", Native: "עברית", RTL: true},
	{Code: "fi", Name: "Finnish", Native: "Suomi"},
	{Code: "da", Name: "Danish", Native: "Dansk"},
	{Code: "cs", Name: "Czech", Native: "Čeština"},
	{Code: "ro", Name: "Romanian", Native: "Română"},
	{Code: "hu", Name: "Hungarian", Native: "Magyar"},
	{Code: "cy", Name: "Welsh", Native: "Cymraeg"},
	{Code: "sw", Name: "Swahili", Native: "Kiswahili"},
	{Code: "id", Name: "Indonesian", Native: "Bahasa Indonesia"},
	{Code: "th", Name: "Thai", Native: "ไทย"},
	{Code: "eo", Name: "Esperanto", Native: "Esperanto"},
	{Code: "la", Name: "Latin", Native: "Latina"},
}

// byCode indexes Languages once at init.
var byCode = func() map[string]Language {
	m := make(map[string]Language, len(Languages))
	for _, l := range Languages {
		m[l.Code] = l
	}
	return m
}()

// Lookup returns the language for an ISO code.
func Lookup(code string) (Language, bool) {
	l, ok := byCode[code]
	return l, ok
}

// ID builds a course id from a base and target code.
func ID(base, target string) string {
	return fmt.Sprintf("%s-from-%s", target, base)
}

// SeedCatalog lists every seed language as a target from the given base.
func SeedCatalog(base string) []Course {
	b, ok := byCode[base]
	if !ok {
		return nil
	}
	out := make([]Course, 0, len(Languages)-1)
	for _, t := range Languages {
		if t.Code == base {
			continue
		}
		out = append(out, Course{ID: ID(base, t.Code), Base: b, Target: t})
	}
	return out
}
