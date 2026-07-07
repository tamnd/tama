package course

import (
	"context"
	"encoding/json"
	"fmt"
)

// Course statuses. Seed courses without a cached pack sit at StatusSeed;
// the long-tail request flow moves a course through generating to ready.
const (
	StatusSeed       = "seed"
	StatusGenerating = "generating"
	StatusReady      = "ready"
	StatusFailed     = "failed"
)

// CEFRSpan is the band a course covers, "A1" to "B2".
type CEFRSpan struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Course is the course row the API serves. Title is in the base language
// ("Japanese" when base is en). Timestamps are unix milliseconds UTC.
type Course struct {
	ID           string   `json:"id"`
	TargetLang   string   `json:"targetLang"`
	BaseLang     string   `json:"baseLang"`
	Title        string   `json:"title"`
	CEFRSpan     CEFRSpan `json:"cefrSpan"`
	Status       string   `json:"status"` // seed, ready, generating, failed
	PackVersion  string   `json:"packVersion"`
	GeneratorVer string   `json:"generatorVersion"`
	CreatedAt    int64    `json:"createdAt"`
}

// Guidebook is the tips page behind the notebook icon on a unit banner,
// always readable regardless of node progress.
type Guidebook struct {
	ID          string       `json:"id"`
	UnitID      string       `json:"unitId"`
	Title       string       `json:"title"` // the unit theme
	KeyPhrases  []KeyPhrase  `json:"keyPhrases"`
	TipSections []TipSection `json:"tipSections"`
}

// KeyPhrase is one card in the guidebook's key phrase list; the speaker
// button plays AudioRef.
type KeyPhrase struct {
	TargetText string `json:"targetText"`
	BaseText   string `json:"baseText"`
	AudioRef   string `json:"audioRef"`
}

// TipSection is one explanatory block; body markdown is restricted to bold,
// italic, and tables, validated at pack load.
type TipSection struct {
	Heading      string    `json:"heading"`
	BodyMarkdown string    `json:"bodyMarkdown"`
	Examples     []Example `json:"examples,omitempty"`
}

// Example pairs a target-language snippet with its base-language reading.
type Example struct {
	Target string `json:"target"`
	Base   string `json:"base"`
}

// Story is an interactive dialogue behind a book node.
type Story struct {
	ID         string          `json:"id"`
	UnitID     string          `json:"unitId"`
	Title      string          `json:"title"`     // target language
	TitleBase  string          `json:"titleBase"` // base language
	CEFR       string          `json:"cefr"`
	Characters []Character     `json:"characters"`
	Lines      []StoryLine     `json:"lines"`
	Exercises  []StoryExercise `json:"exercises,omitempty"`
	XPReward   int             `json:"xpReward"` // default 20, awarded once
}

// Character is one recurring cast member; the cast generates once per course
// and reappears across stories so learners recognize it.
type Character struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarRef string `json:"avatarRef"`
	Voice     string `json:"voice"`
}

// StoryLine is one dialogue line; SpeakerID is empty for narration.
type StoryLine struct {
	SpeakerID  string      `json:"speakerId"`
	TargetText string      `json:"targetText"`
	BaseText   string      `json:"baseText"`
	AudioRef   string      `json:"audioRef,omitempty"`
	TokenHints []TokenHint `json:"tokenHints,omitempty"`
}

// TokenHint carries per-token hover data; Reading feeds furigana and pinyin
// rendering for CJK targets.
type TokenHint struct {
	Token   string `json:"token"`
	Hint    string `json:"hint,omitempty"`
	Reading string `json:"reading,omitempty"`
}

// StoryExercise interleaves a comprehension exercise after a line. The
// exercise body reuses the M4 exercise envelope and stays opaque here.
type StoryExercise struct {
	AfterLine int             `json:"afterLine"` // 1-based line index
	Exercise  json.RawMessage `json:"exercise"`
}

// Validate checks the structural story rules: a cast, lines whose speakers
// exist, and embedded exercises anchored to real lines.
func (s Story) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("story: missing id")
	}
	if s.Title == "" || s.TitleBase == "" {
		return fmt.Errorf("story %s: missing title", s.ID)
	}
	if len(s.Characters) == 0 {
		return fmt.Errorf("story %s: no characters", s.ID)
	}
	if len(s.Lines) == 0 {
		return fmt.Errorf("story %s: no lines", s.ID)
	}
	cast := make(map[string]bool, len(s.Characters))
	for _, c := range s.Characters {
		if c.ID == "" || c.Name == "" {
			return fmt.Errorf("story %s: character missing id or name", s.ID)
		}
		if cast[c.ID] {
			return fmt.Errorf("story %s: duplicate character %s", s.ID, c.ID)
		}
		cast[c.ID] = true
	}
	for i, l := range s.Lines {
		if l.SpeakerID != "" && !cast[l.SpeakerID] {
			return fmt.Errorf("story %s: line %d speaker %q not in cast", s.ID, i+1, l.SpeakerID)
		}
		if l.TargetText == "" {
			return fmt.Errorf("story %s: line %d has no target text", s.ID, i+1)
		}
	}
	for _, e := range s.Exercises {
		if e.AfterLine < 1 || e.AfterLine > len(s.Lines) {
			return fmt.Errorf("story %s: exercise after line %d, story has %d lines", s.ID, e.AfterLine, len(s.Lines))
		}
	}
	return nil
}

// Path is the single response that renders the whole path screen: ordered
// sections, each with ordered units, each with ordered nodes.
type Path struct {
	CourseID string        `json:"courseId"`
	Sections []PathSection `json:"sections"`
}

// PathSection is one section band in the path response.
type PathSection struct {
	Index int        `json:"index"`
	Title string     `json:"title"` // "Section 1: Rookie"
	CEFR  string     `json:"cefr"`
	Color string     `json:"color"`
	Units []PathUnit `json:"units"`
	Gate  *PathGate  `json:"gate,omitempty"`
}

// PathGate is the jump-ahead entry at a section boundary.
type PathGate struct {
	NodeID string `json:"nodeId"`
	Label  string `json:"label"` // "JUMP HERE?"
}

// PathUnit is one unit banner plus its ordered nodes.
type PathUnit struct {
	ID          string `json:"id"`
	Index       int    `json:"index"`
	Theme       string `json:"theme"`
	GuidebookID string `json:"guidebookId,omitempty"`
	Color       string `json:"color"`
	Nodes       []Node `json:"nodes"`
}

// Progress is the per-course rollup the profile screen reads.
type Progress struct {
	Crowns         int `json:"crowns"`
	LegendaryCount int `json:"legendaryCount"`
	NodesDone      int `json:"nodesDone"`
	TotalNodes     int `json:"totalNodes"`
}

// Store is the read model the API layer works against. Pack import compiles
// pack JSON into SQLite and these methods read only SQLite, never pack
// files; pkg/store implements it, and no query strings live outside there.
type Store interface {
	GetCourse(ctx context.Context, id string) (Course, error)
	ListCourses(ctx context.Context) ([]Course, error)
	// GetPath returns the full node list for a course with the user's
	// states resolved.
	GetPath(ctx context.Context, courseID string, userID int64) (Path, error)
	GetGuidebook(ctx context.Context, unitID string) (Guidebook, error)
	GetStory(ctx context.Context, storyID string) (Story, error)
	GetLesson(ctx context.Context, lessonID string) (Lesson, error)
	// InstallPack validates the pack directory and imports it
	// transactionally; a failed import leaves the installed version alone.
	InstallPack(ctx context.Context, dir string) error
}
