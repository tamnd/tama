package course

import "testing"

func TestNodeKindValid(t *testing.T) {
	for _, k := range []NodeKind{KindLesson, KindChest, KindHard, KindStory, KindReview, KindCheckpoint, KindGate, KindPractice} {
		if !k.Valid() {
			t.Errorf("kind %q should be valid", k)
		}
	}
	if NodeKind("boss").Valid() {
		t.Error("kind boss should be invalid")
	}
}

func TestSectionHelpers(t *testing.T) {
	if got := SectionName(1); got != "Rookie" {
		t.Errorf("SectionName(1) = %q", got)
	}
	if got := SectionName(5); got != "Champion" {
		t.Errorf("SectionName(5) = %q", got)
	}
	if got := SectionName(8); got != "Champion" {
		t.Errorf("SectionName(8) = %q, long courses keep the last name", got)
	}
	if got := SectionTitle(1, "Rookie"); got != "Section 1: Rookie" {
		t.Errorf("SectionTitle = %q", got)
	}
	if got := SectionColor(1); got != "green" {
		t.Errorf("SectionColor(1) = %q", got)
	}
	if got := SectionColor(7); got != "green" {
		t.Errorf("SectionColor(7) = %q, palette cycles", got)
	}
}

func TestSectionValidate(t *testing.T) {
	good := Section{ID: "s1", CourseID: "es-from-en", Index: 1, Title: "Section 1: Rookie", CEFRLabel: "A1", Color: "green", UnitCount: 4}
	if err := good.Validate(); err != nil {
		t.Errorf("valid section: %v", err)
	}
	tests := []struct {
		name string
		mut  func(*Section)
	}{
		{"missing id", func(s *Section) { s.ID = "" }},
		{"zero index", func(s *Section) { s.Index = 0 }},
		{"bad cefr", func(s *Section) { s.CEFRLabel = "Z9" }},
		{"bad color", func(s *Section) { s.Color = "chartreuse" }},
		{"no units", func(s *Section) { s.UnitCount = 0 }},
	}
	for _, tt := range tests {
		s := good
		tt.mut(&s)
		if err := s.Validate(); err == nil {
			t.Errorf("%s: want error", tt.name)
		}
	}
}

func TestUnitValidate(t *testing.T) {
	good := Unit{ID: "u1", SectionIndex: 1, Index: 1, Theme: "Order food and drink", Title: "Unit 1", Color: "green", GuidebookID: "gb_u1", LevelCount: 4}
	if err := good.Validate(); err != nil {
		t.Errorf("valid unit: %v", err)
	}
	tests := []struct {
		name string
		mut  func(*Unit)
	}{
		{"missing id", func(u *Unit) { u.ID = "" }},
		{"zero section", func(u *Unit) { u.SectionIndex = 0 }},
		{"no theme", func(u *Unit) { u.Theme = "" }},
		{"bad color", func(u *Unit) { u.Color = "mauve" }},
		{"no levels", func(u *Unit) { u.LevelCount = 0 }},
	}
	for _, tt := range tests {
		u := good
		tt.mut(&u)
		if err := u.Validate(); err == nil {
			t.Errorf("%s: want error", tt.name)
		}
	}
}

func TestLevelValidate(t *testing.T) {
	tests := []struct {
		name string
		l    Level
		ok   bool
	}{
		{"regular", Level{ID: "n1", UnitID: "u1", Index: 1, Kind: KindLesson, LessonCount: 4, HardLessonCount: 1}, true},
		{"regular six", Level{ID: "n1", UnitID: "u1", Index: 1, Kind: KindLesson, LessonCount: 6}, true},
		{"regular three", Level{ID: "n1", UnitID: "u1", Index: 1, Kind: KindLesson, LessonCount: 3}, false},
		{"regular seven", Level{ID: "n1", UnitID: "u1", Index: 1, Kind: KindLesson, LessonCount: 7}, false},
		{"hard", Level{ID: "n2", UnitID: "u1", Index: 2, Kind: KindHard, LessonCount: 1}, true},
		{"hard two", Level{ID: "n2", UnitID: "u1", Index: 2, Kind: KindHard, LessonCount: 2}, false},
		{"chest", Level{ID: "n3", UnitID: "u1", Index: 3, Kind: KindChest}, true},
		{"chest with lessons", Level{ID: "n3", UnitID: "u1", Index: 3, Kind: KindChest, LessonCount: 1}, false},
		{"story", Level{ID: "n4", UnitID: "u1", Index: 4, Kind: KindStory}, true},
		{"review", Level{ID: "n5", UnitID: "u1", Index: 5, Kind: KindReview, LessonCount: 1}, true},
		{"checkpoint", Level{ID: "n6", UnitID: "u1", Index: 6, Kind: KindCheckpoint, LessonCount: 1}, true},
		{"gate", Level{ID: "n7", UnitID: "u1", Index: 7, Kind: KindGate}, true},
		{"practice", Level{ID: "n8", UnitID: "u1", Index: 8, Kind: KindPractice}, true},
		{"bad kind", Level{ID: "n9", UnitID: "u1", Index: 9, Kind: "boss"}, false},
		{"hard count high", Level{ID: "n1", UnitID: "u1", Index: 1, Kind: KindLesson, LessonCount: 4, HardLessonCount: 5}, false},
		{"no unit", Level{ID: "n1", Index: 1, Kind: KindLesson, LessonCount: 4}, false},
	}
	for _, tt := range tests {
		if err := tt.l.Validate(); (err == nil) != tt.ok {
			t.Errorf("%s: Validate() = %v, want ok=%v", tt.name, err, tt.ok)
		}
	}
}

func TestLessonValidate(t *testing.T) {
	good := Lesson{ID: "ls1", LevelID: "n1", Index: 1, ExerciseIDs: []string{"ex1"}, XPReward: XPLesson}
	if err := good.Validate(); err != nil {
		t.Errorf("valid lesson: %v", err)
	}
	for name, mut := range map[string]func(*Lesson){
		"missing id":    func(l *Lesson) { l.ID = "" },
		"missing level": func(l *Lesson) { l.LevelID = "" },
		"zero index":    func(l *Lesson) { l.Index = 0 },
		"zero xp":       func(l *Lesson) { l.XPReward = 0 },
	} {
		l := good
		mut(&l)
		if err := l.Validate(); err == nil {
			t.Errorf("%s: want error", name)
		}
	}
}

// TestDeriveStates is the icon-state table: locked, active, completed, and
// legendary, with exactly one active node until the path is done.
func TestDeriveStates(t *testing.T) {
	done := NodeProgress{Completed: true}
	leg := NodeProgress{Completed: true, Legendary: true}
	fresh := NodeProgress{}
	started := NodeProgress{LessonsDone: 2}

	tests := []struct {
		name string
		prog []NodeProgress
		want []NodeState
	}{
		{"fresh course", []NodeProgress{fresh, fresh, fresh},
			[]NodeState{StateActive, StateLocked, StateLocked}},
		{"mid course", []NodeProgress{done, started, fresh},
			[]NodeState{StateCompleted, StateActive, StateLocked}},
		{"legendary node", []NodeProgress{leg, done, fresh},
			[]NodeState{StateLegendary, StateCompleted, StateActive}},
		{"all done", []NodeProgress{done, done, leg},
			[]NodeState{StateCompleted, StateCompleted, StateLegendary}},
		{"jump gate skipped ahead", []NodeProgress{done, fresh, done, fresh},
			[]NodeState{StateCompleted, StateActive, StateCompleted, StateLocked}},
		{"empty path", nil, []NodeState{}},
	}
	for _, tt := range tests {
		got := DeriveStates(tt.prog)
		if len(got) != len(tt.want) {
			t.Errorf("%s: %d states, want %d", tt.name, len(got), len(tt.want))
			continue
		}
		active := 0
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("%s: state[%d] = %s, want %s", tt.name, i, got[i], tt.want[i])
			}
			if got[i] == StateActive {
				active++
			}
		}
		if active > 1 {
			t.Errorf("%s: %d active nodes, want at most one", tt.name, active)
		}
	}
}

func TestStoryValidate(t *testing.T) {
	good := Story{
		ID: "st1", UnitID: "u1", Title: "El café de Tama", TitleBase: "Tama's café", CEFR: "A1",
		Characters: []Character{{ID: "tama", Name: "Tama", AvatarRef: "cast/tama.svg", Voice: "f1"}},
		Lines: []StoryLine{
			{SpeakerID: "tama", TargetText: "Un café, por favor.", BaseText: "A coffee, please."},
			{SpeakerID: "", TargetText: "Tama espera.", BaseText: "Tama waits."},
		},
		Exercises: []StoryExercise{{AfterLine: 2, Exercise: []byte(`{}`)}},
		XPReward:  XPStory,
	}
	if err := good.Validate(); err != nil {
		t.Errorf("valid story: %v", err)
	}
	tests := []struct {
		name string
		mut  func(*Story)
	}{
		{"no cast", func(s *Story) { s.Characters = nil }},
		{"no lines", func(s *Story) { s.Lines = nil }},
		{"unknown speaker", func(s *Story) { s.Lines[0].SpeakerID = "mesero" }},
		{"exercise off the end", func(s *Story) { s.Exercises[0].AfterLine = 9 }},
		{"exercise before start", func(s *Story) { s.Exercises[0].AfterLine = 0 }},
		{"no base title", func(s *Story) { s.TitleBase = "" }},
	}
	for _, tt := range tests {
		s := good
		s.Characters = append([]Character(nil), good.Characters...)
		s.Lines = append([]StoryLine(nil), good.Lines...)
		s.Exercises = append([]StoryExercise(nil), good.Exercises...)
		tt.mut(&s)
		if err := s.Validate(); err == nil {
			t.Errorf("%s: want error", tt.name)
		}
	}
}
