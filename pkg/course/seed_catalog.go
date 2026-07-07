package course

import "github.com/tamnd/tama/pkg/course/lang"

// SeedCourse is one out-of-the-box picker entry. Seed courses are still
// generated through the normal pack pipeline (M8); the seed list marks them
// as first-class picker entries and generation-priority targets.
type SeedCourse struct {
	ID         string `json:"id"`
	TargetLang string `json:"targetLang"` // shortest code, e.g. "es"
	BaseLang   string `json:"baseLang"`   // shortest code, e.g. "en"
	// Title is the course title the picker shows. Seed titles use the
	// target's English name; per-base localized titles arrive with the
	// generated packs.
	Title    string `json:"title"`
	Order    int    `json:"order"` // display order within its base group
	Flag     string `json:"flag"`  // flag asset ref, e.g. "flags/es.svg"
	Learners string `json:"learners"`
	// Description carries course-level notes, like Portuguese teaching the
	// Brazilian variant.
	Description string `json:"description,omitempty"`
	// ScriptOptions lists selectable target scripts where the course has a
	// setting for it; Chinese defaults to Hans with a Hant toggle.
	ScriptOptions []string `json:"scriptOptions,omitempty"`
}

// seedTarget pairs a target code with its cosmetic learner count.
type seedTarget struct {
	code, learners string
}

// fromEnglish is the from-English seed list in the picker's display order,
// mirroring the real product's popularity ladder: Spanish first, French
// second, down through the constructed and fictional tail.
var fromEnglish = []seedTarget{
	{"es", "48.4M"}, {"fr", "27.3M"}, {"ja", "17.9M"}, {"de", "14.5M"},
	{"ko", "14.2M"}, {"it", "10.4M"}, {"zh", "8.99M"}, {"ru", "7.46M"},
	{"hi", "7.14M"}, {"pt", "5.31M"}, {"ar", "5.05M"}, {"nl", "3.35M"},
	{"sv", "2.62M"}, {"el", "2.24M"}, {"ga", "1.99M"}, {"pl", "1.85M"},
	{"nb", "1.79M"}, {"tr", "1.72M"}, {"vi", "1.65M"}, {"uk", "1.61M"},
	{"he", "1.35M"}, {"fi", "1.31M"}, {"da", "1.14M"}, {"cs", "1.06M"},
	{"ro", "927K"}, {"hu", "769K"}, {"cy", "692K"}, {"sw", "634K"},
	{"id", "588K"}, {"haw", "541K"}, {"nv", "296K"}, {"gd", "483K"},
	{"yi", "312K"}, {"ht", "327K"}, {"zu", "377K"}, {"eo", "279K"},
	{"la", "1.02M"}, {"tlh", "268K"}, {"qhv", "631K"}, {"ca", "310K"},
}

// seedBases lists the non-English base languages with the target set the
// real product offers from each, English always first.
var seedBases = []struct {
	base    string
	targets []seedTarget
}{
	{"es", []seedTarget{{"en", "35.9M"}, {"fr", "3.1M"}, {"pt", "2.6M"}, {"it", "2.2M"}, {"de", "1.9M"}, {"ca", "410K"}, {"gn", "230K"}}},
	{"pt", []seedTarget{{"en", "12.4M"}, {"es", "2.9M"}, {"fr", "1.6M"}, {"de", "890K"}, {"it", "740K"}, {"eo", "150K"}}},
	{"zh", []seedTarget{{"en", "10.2M"}, {"ja", "2.4M"}, {"ko", "1.7M"}, {"fr", "820K"}, {"es", "760K"}, {"it", "430K"}}},
	{"ja", []seedTarget{{"en", "6.1M"}, {"zh", "1.3M"}, {"ko", "1.1M"}, {"fr", "460K"}}},
	{"ko", []seedTarget{{"en", "4.5M"}, {"ja", "1.5M"}, {"zh", "870K"}}},
	{"fr", []seedTarget{{"en", "8.3M"}, {"es", "2.4M"}, {"de", "1.3M"}, {"it", "1.2M"}, {"pt", "580K"}, {"eo", "170K"}}},
	{"de", []seedTarget{{"en", "3.2M"}, {"es", "1.4M"}, {"fr", "1.1M"}}},
	{"ru", []seedTarget{{"en", "5.6M"}, {"es", "980K"}, {"de", "870K"}, {"fr", "640K"}}},
	{"hi", []seedTarget{{"en", "9.8M"}}},
	{"ar", []seedTarget{{"en", "6.4M"}, {"fr", "1.2M"}, {"de", "540K"}, {"sv", "300K"}}},
	{"vi", []seedTarget{{"en", "4.9M"}, {"zh", "720K"}, {"ja", "680K"}, {"ko", "610K"}}},
	{"tr", []seedTarget{{"en", "3.7M"}, {"de", "740K"}, {"ru", "350K"}}},
	{"it", []seedTarget{{"en", "2.8M"}, {"es", "910K"}, {"fr", "830K"}, {"de", "560K"}}},
	{"id", []seedTarget{{"en", "3.4M"}}},
	{"th", []seedTarget{{"en", "2.1M"}}},
	{"uk", []seedTarget{{"en", "1.9M"}, {"es", "210K"}}},
	{"pl", []seedTarget{{"en", "1.4M"}}},
	{"cs", []seedTarget{{"en", "820K"}}},
	{"hu", []seedTarget{{"en", "610K"}}},
	{"el", []seedTarget{{"en", "540K"}}},
	{"bn", []seedTarget{{"en", "2.6M"}}},
	{"te", []seedTarget{{"en", "1.8M"}}},
	{"ta", []seedTarget{{"en", "1.7M"}}},
	{"nl", []seedTarget{{"en", "930K"}}},
}

// seedCourses is built once from the tables above.
var seedCourses = buildSeedCourses()

func buildSeedCourses() []SeedCourse {
	var out []SeedCourse
	add := func(baseCode string, targets []seedTarget) {
		base, err := lang.Lookup(baseCode)
		if err != nil {
			panic(err)
		}
		for i, t := range targets {
			target, err := lang.Lookup(t.code)
			if err != nil {
				panic(err)
			}
			c := SeedCourse{
				ID:         ID(target, base),
				TargetLang: target.Code(),
				BaseLang:   base.Code(),
				Title:      target.Name,
				Order:      i + 1,
				Flag:       "flags/" + target.Code() + ".svg",
				Learners:   t.learners,
			}
			switch target.Code() {
			case "pt":
				c.Description = "Teaches Brazilian Portuguese."
			case "zh":
				c.Description = "Teaches Mandarin, simplified characters by default; a course setting switches to traditional."
				c.ScriptOptions = []string{"Hans", "Hant"}
			}
			out = append(out, c)
		}
	}
	add("en", fromEnglish)
	for _, b := range seedBases {
		add(b.base, b.targets)
	}
	return out
}

// SeedCourses returns every seed course across all bases. The slice is a
// copy; callers may reorder it.
func SeedCourses() []SeedCourse {
	out := make([]SeedCourse, len(seedCourses))
	copy(out, seedCourses)
	return out
}

// SeedCatalog returns the seed picker entries for one base language in
// display order, nil when the base has no seed courses. Any valid pair
// outside the seed list stays addressable through the long-tail request
// flow; this is only what the picker surfaces first.
func SeedCatalog(base string) []SeedCourse {
	b, err := lang.Lookup(base)
	if err != nil {
		return nil
	}
	var out []SeedCourse
	for _, c := range seedCourses {
		if c.BaseLang == b.Code() {
			out = append(out, c)
		}
	}
	return out
}
