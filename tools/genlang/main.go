// Command genlang turns the ISO 639-3 code table into the compiled-in
// language registry at pkg/course/lang/langdata.go. It runs once per table
// update via go generate; the server never fetches language data at runtime.
//
// The input is iso-639-3.tsv next to this file, a slim projection of the
// ISO 639-3 registry: alpha3, alpha2, English name, scope, type. On top of
// that sit curated overrides for the languages the seed catalog touches:
// endonyms, scripts, TTS coverage, and the private-use and constructed rows
// (High Valyrian, Klingon) the picker needs.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
)

// override carries the curated fields for one 639-3 code.
type override struct {
	part1  string // set or replace the 639-1 code
	name   string // replace the English name shown in the picker
	native string // endonym
	script string // ISO 15924
	tts    bool   // the audio pipeline has a voice
}

// overrides is the curated layer: every language the seed catalog references
// gets its endonym and script here, plus the handful of extra rows the
// script-aware milestones need (fa, ur, dv, syr for RTL coverage).
var overrides = map[string]override{
	"eng": {native: "English", script: "Latn", tts: true},
	"spa": {native: "Español", script: "Latn", tts: true},
	"fra": {native: "Français", script: "Latn", tts: true},
	"jpn": {native: "日本語", script: "Jpan", tts: true},
	"deu": {native: "Deutsch", script: "Latn", tts: true},
	"kor": {native: "한국어", script: "Hang", tts: true},
	"ita": {native: "Italiano", script: "Latn", tts: true},
	// Chinese folds onto Mandarin: zh, zho, and cmn all reach this row, and
	// courses teach simplified by default (zh-Hans) with a traditional
	// setting on the course.
	"cmn": {part1: "zh", name: "Chinese", native: "中文", script: "Hani", tts: true},
	"rus": {native: "Русский", script: "Cyrl", tts: true},
	"hin": {native: "हिन्दी", script: "Deva", tts: true},
	"por": {native: "Português", script: "Latn", tts: true},
	"ara": {native: "العربية", script: "Arab", tts: true},
	"nld": {native: "Nederlands", script: "Latn", tts: true},
	"swe": {native: "Svenska", script: "Latn", tts: true},
	"ell": {name: "Greek", native: "Ελληνικά", script: "Grek", tts: true},
	"gle": {native: "Gaeilge", script: "Latn", tts: true},
	"pol": {native: "Polski", script: "Latn", tts: true},
	// The product teaches Bokmål and labels it that way in the picker.
	"nob": {name: "Norwegian (Bokmål)", native: "Norsk (bokmål)", script: "Latn", tts: true},
	"nor": {native: "Norsk", script: "Latn", tts: true},
	"tur": {native: "Türkçe", script: "Latn", tts: true},
	"vie": {native: "Tiếng Việt", script: "Latn", tts: true},
	"ukr": {native: "Українська", script: "Cyrl", tts: true},
	"heb": {native: "עברית", script: "Hebr", tts: true},
	"fin": {native: "Suomi", script: "Latn", tts: true},
	"dan": {native: "Dansk", script: "Latn", tts: true},
	"ces": {native: "Čeština", script: "Latn", tts: true},
	"ron": {native: "Română", script: "Latn", tts: true},
	"hun": {native: "Magyar", script: "Latn", tts: true},
	"cym": {native: "Cymraeg", script: "Latn", tts: true},
	"swh": {name: "Swahili (coastal)", native: "Kiswahili", script: "Latn", tts: true},
	"swa": {name: "Swahili", native: "Kiswahili", script: "Latn", tts: true},
	"ind": {native: "Bahasa Indonesia", script: "Latn", tts: true},
	"haw": {native: "ʻŌlelo Hawaiʻi", script: "Latn", tts: true},
	"nav": {native: "Diné bizaad", script: "Latn"},
	"gla": {native: "Gàidhlig", script: "Latn", tts: true},
	"yid": {native: "ייִדיש", script: "Hebr", tts: true},
	"hat": {name: "Haitian Creole", native: "Kreyòl ayisyen", script: "Latn", tts: true},
	"zul": {native: "isiZulu", script: "Latn", tts: true},
	"epo": {native: "Esperanto", script: "Latn", tts: true},
	"lat": {native: "Latina", script: "Latn"},
	"tlh": {native: "tlhIngan Hol", script: "Latn"},
	"cat": {native: "Català", script: "Latn", tts: true},
	"tha": {native: "ไทย", script: "Thai", tts: true},
	"ben": {native: "বাংলা", script: "Beng", tts: true},
	"tel": {native: "తెలుగు", script: "Telu", tts: true},
	"tam": {native: "தமிழ்", script: "Taml", tts: true},
	"grn": {native: "Avañe'ẽ", script: "Latn"},
	"fas": {native: "فارسی", script: "Arab", tts: true},
	"urd": {native: "اردو", script: "Arab", tts: true},
	"div": {native: "ދިވެހި", script: "Thaa"},
	"syr": {native: "ܣܘܪܝܝܐ", script: "Syrc"},
}

// extraRows are codes outside ISO 639-3: private-use fictional languages.
// qaa..qtz is the 639-3 private-use range.
var extraRows = []row{
	{alpha3: "qhv", name: "High Valyrian", native: "Valyrio", script: "Latn"},
}

// dropRows removes table rows whose codes become aliases of another row.
// zho (Chinese, macrolanguage) folds onto cmn so zh reaches Mandarin.
var dropRows = map[string]bool{"zho": true}

// aliases maps extra lookup codes onto the 639-3 code that owns the row.
var aliases = map[string]string{
	"zho": "cmn",
}

// rtlScripts drives Direction; exactly these scripts read right to left in
// the seed data.
var rtlScripts = map[string]bool{"Arab": true, "Hebr": true, "Syrc": true, "Thaa": true}

type row struct {
	alpha3, alpha1, name, native, script string
	tts                                  bool
}

func main() {
	tsv := flag.String("tsv", "iso-639-3.tsv", "path to the ISO 639-3 table")
	out := flag.String("o", "langdata.go", "output file")
	flag.Parse()

	raw, err := os.ReadFile(*tsv)
	if err != nil {
		log.Fatal(err)
	}

	var rows []row
	for i, line := range strings.Split(string(raw), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 5 {
			log.Fatalf("%s:%d: want 5 fields, got %d", *tsv, i+1, len(f))
		}
		alpha3, alpha1, name, typ := f[0], f[1], f[2], f[4]
		if typ == "S" || dropRows[alpha3] { // special codes (mis, mul, und, zxx)
			continue
		}
		r := row{alpha3: alpha3, alpha1: alpha1, name: name}
		if o, ok := overrides[alpha3]; ok {
			if o.part1 != "" {
				r.alpha1 = o.part1
			}
			if o.name != "" {
				r.name = o.name
			}
			r.native, r.script, r.tts = o.native, o.script, o.tts
		}
		rows = append(rows, r)
	}
	rows = append(rows, extraRows...)
	sort.Slice(rows, func(i, j int) bool { return rows[i].alpha3 < rows[j].alpha3 })

	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by tools/genlang; DO NOT EDIT.\n\npackage lang\n\n")
	fmt.Fprintf(&b, "// languages is every registry row, sorted by 639-3 code. Empty\n")
	fmt.Fprintf(&b, "// NativeName defaults to Name and empty Direction to ltr at init.\n")
	fmt.Fprintf(&b, "var languages = []Language{\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "\t{")
		if r.alpha1 != "" {
			fmt.Fprintf(&b, "Code639_1: %q, ", r.alpha1)
		}
		fmt.Fprintf(&b, "Code639_3: %q, Name: %q", r.alpha3, r.name)
		if r.native != "" && r.native != r.name {
			fmt.Fprintf(&b, ", NativeName: %q", r.native)
		}
		if r.script != "" {
			fmt.Fprintf(&b, ", Script: %q", r.script)
		}
		if rtlScripts[r.script] {
			fmt.Fprintf(&b, ", Direction: DirectionRTL")
		}
		if r.tts {
			fmt.Fprintf(&b, ", HasTTS: true")
		}
		fmt.Fprintf(&b, "},\n")
	}
	fmt.Fprintf(&b, "}\n\n")

	fmt.Fprintf(&b, "// aliases maps extra lookup codes to the 639-3 code owning the row.\n")
	fmt.Fprintf(&b, "var aliases = map[string]string{\n")
	keys := make([]string, 0, len(aliases))
	for k := range aliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&b, "\t%q: %q,\n", k, aliases[k])
	}
	fmt.Fprintf(&b, "}\n")

	src, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatalf("format: %v", err)
	}
	if err := os.WriteFile(*out, src, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s: %d languages\n", *out, len(rows))
}
