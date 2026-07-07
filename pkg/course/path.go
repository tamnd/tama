package course

import (
	"encoding/json"
	"fmt"
)

// XP rewards. A hard node's single lesson pays double, story completion pays
// its own reward once, and replaying a finished lesson pays the reduced rate.
const (
	XPLesson       = 10
	XPLessonReplay = 5
	XPHard         = 2 * XPLesson
	XPStory        = 20
)

// NodeKind is the discriminator every path node carries; the renderer
// switches icon, color, and tap behavior on it.
type NodeKind string

const (
	KindLesson     NodeKind = "lesson"     // regular star node
	KindChest      NodeKind = "chest"      // reward node, grants gems once
	KindHard       NodeKind = "hard"       // single harder lesson, double XP
	KindStory      NodeKind = "story"      // book node linking to a story
	KindReview     NodeKind = "review"     // unit review, legendary entry point
	KindCheckpoint NodeKind = "checkpoint" // unit trophy at the end of every unit
	KindGate       NodeKind = "gate"       // jump-ahead test-out at unit boundaries
	KindPractice   NodeKind = "practice"   // personalized practice, injected by M6
)

// nodeKinds is the closed set.
var nodeKinds = map[NodeKind]bool{
	KindLesson: true, KindChest: true, KindHard: true, KindStory: true,
	KindReview: true, KindCheckpoint: true, KindGate: true, KindPractice: true,
}

// Valid reports whether k is one of the eight node kinds.
func (k NodeKind) Valid() bool { return nodeKinds[k] }

// NodeState is the derived icon state of a path node for one user.
type NodeState string

const (
	StateLocked    NodeState = "locked"    // gray
	StateActive    NodeState = "active"    // section color, bouncing START pin
	StateCompleted NodeState = "completed" // gold ring
	StateLegendary NodeState = "legendary" // purple with wreath
)

// Valid reports whether s is one of the four node states.
func (s NodeState) Valid() bool {
	switch s {
	case StateLocked, StateActive, StateCompleted, StateLegendary:
		return true
	}
	return false
}

// CEFR labels a section may carry. Section 1 is always pre-A1/A1 territory
// and labels itself A1.
var cefrLabels = map[string]bool{"A1": true, "A2": true, "B1": true, "B2": true, "C1": true}

// SectionColors are the section palette roles from M2, in the order sections
// cycle through them. They map onto the token palette: green is
// feather-green, blue is macaw, purple is beetle, orange is fox, red is
// cardinal, yellow is bee.
var SectionColors = []string{"green", "blue", "purple", "orange", "red", "yellow"}

// sectionColorSet indexes SectionColors for validation.
var sectionColorSet = func() map[string]bool {
	m := make(map[string]bool, len(SectionColors))
	for _, c := range SectionColors {
		m[c] = true
	}
	return m
}()

// SectionColor is the palette role for a 1-based section index; sections
// cycle through the palette in order.
func SectionColor(index int) string {
	return SectionColors[(index-1)%len(SectionColors)]
}

// sectionNames is the product's section ladder. Courses longer than the
// ladder keep the last name.
var sectionNames = []string{"Rookie", "Explorer", "Traveler", "Trailblazer", "Champion"}

// SectionName is the ladder name for a 1-based section index.
func SectionName(index int) string {
	if index <= len(sectionNames) {
		return sectionNames[index-1]
	}
	return sectionNames[len(sectionNames)-1]
}

// SectionTitle formats the section header title, "Section 1: Rookie".
func SectionTitle(index int, name string) string {
	return fmt.Sprintf("Section %d: %s", index, name)
}

// Section is one CEFR band of a course. Courses have 3 to 8 of them.
type Section struct {
	ID        string `json:"id"`
	CourseID  string `json:"courseId"`
	Index     int    `json:"index"` // 1-based
	Title     string `json:"title"` // "Section 1: Rookie"
	CEFRLabel string `json:"cefrLabel"`
	Color     string `json:"color"` // one of SectionColors
	UnitCount int    `json:"unitCount"`
}

// Validate checks the structural rules for a section row.
func (s Section) Validate() error {
	switch {
	case s.ID == "":
		return fmt.Errorf("section: missing id")
	case s.Index < 1:
		return fmt.Errorf("section %s: index %d, want >= 1", s.ID, s.Index)
	case !cefrLabels[s.CEFRLabel]:
		return fmt.Errorf("section %s: bad cefr label %q", s.ID, s.CEFRLabel)
	case !sectionColorSet[s.Color]:
		return fmt.Errorf("section %s: color %q is not a section palette role", s.ID, s.Color)
	case s.UnitCount < 1:
		return fmt.Errorf("section %s: no units", s.ID)
	}
	return nil
}

// Unit is one themed band of a section. The banner shows the theme, the
// unit number, and the guidebook button.
type Unit struct {
	ID           string `json:"id"`
	SectionIndex int    `json:"sectionIndex"` // 1-based section it belongs to
	Index        int    `json:"index"`        // 1-based within the section
	Theme        string `json:"theme"`        // "Order food and drink"
	Title        string `json:"title"`        // "Unit 12"
	Color        string `json:"color"`        // inherited from the section
	GuidebookID  string `json:"guidebookId,omitempty"`
	LevelCount   int    `json:"levelCount"`
}

// Validate checks the structural rules for a unit row.
func (u Unit) Validate() error {
	switch {
	case u.ID == "":
		return fmt.Errorf("unit: missing id")
	case u.SectionIndex < 1:
		return fmt.Errorf("unit %s: section index %d, want >= 1", u.ID, u.SectionIndex)
	case u.Index < 1:
		return fmt.Errorf("unit %s: index %d, want >= 1", u.ID, u.Index)
	case u.Theme == "":
		return fmt.Errorf("unit %s: missing theme", u.ID)
	case !sectionColorSet[u.Color]:
		return fmt.Errorf("unit %s: color %q is not a section palette role", u.ID, u.Color)
	case u.LevelCount < 1:
		return fmt.Errorf("unit %s: no levels", u.ID)
	}
	return nil
}

// Lesson count bounds per node kind: a regular level carries 4 to 6 lessons,
// single-session kinds carry exactly one, reward and link kinds carry none.
var lessonBounds = map[NodeKind][2]int{
	KindLesson:     {4, 6},
	KindHard:       {1, 1},
	KindReview:     {1, 1},
	KindCheckpoint: {1, 1},
	KindChest:      {0, 0},
	KindStory:      {0, 0},
	KindGate:       {0, 0},
	KindPractice:   {0, 0},
}

// Level is one path node as authored in a pack: position, kind, and payload.
// Node is the same thing after user state resolves on top.
type Level struct {
	ID              string          `json:"id"`
	UnitID          string          `json:"unitId"`
	Index           int             `json:"index"` // 1-based within the unit
	Kind            NodeKind        `json:"kind"`
	LessonCount     int             `json:"lessonCount"`
	HardLessonCount int             `json:"hardLessonCount,omitempty"`
	Payload         json.RawMessage `json:"payload,omitempty"` // kind-discriminated
}

// Validate checks the structural rules for a level row.
func (l Level) Validate() error {
	switch {
	case l.ID == "":
		return fmt.Errorf("level: missing id")
	case l.UnitID == "":
		return fmt.Errorf("level %s: missing unit id", l.ID)
	case l.Index < 1:
		return fmt.Errorf("level %s: index %d, want >= 1", l.ID, l.Index)
	case !l.Kind.Valid():
		return fmt.Errorf("level %s: unknown kind %q", l.ID, l.Kind)
	}
	b := lessonBounds[l.Kind]
	if l.LessonCount < b[0] || l.LessonCount > b[1] {
		return fmt.Errorf("level %s: %s node has %d lessons, want %d to %d",
			l.ID, l.Kind, l.LessonCount, b[0], b[1])
	}
	if l.HardLessonCount < 0 || l.HardLessonCount > l.LessonCount {
		return fmt.Errorf("level %s: hard lesson count %d out of range", l.ID, l.HardLessonCount)
	}
	return nil
}

// Lesson is one session inside a level; lessons within a node are strictly
// sequential.
type Lesson struct {
	ID          string   `json:"id"`
	LevelID     string   `json:"levelId"`
	Index       int      `json:"index"` // 1-based within the level
	ExerciseIDs []string `json:"exerciseIds"`
	XPReward    int      `json:"xpReward"` // base 10
}

// Validate checks the structural rules for a lesson row. Exercise count
// bands are the pack validator's job because one of them is only a warning.
func (l Lesson) Validate() error {
	switch {
	case l.ID == "":
		return fmt.Errorf("lesson: missing id")
	case l.LevelID == "":
		return fmt.Errorf("lesson %s: missing level id", l.ID)
	case l.Index < 1:
		return fmt.Errorf("lesson %s: index %d, want >= 1", l.ID, l.Index)
	case l.XPReward < 1:
		return fmt.Errorf("lesson %s: xp reward %d, want >= 1", l.ID, l.XPReward)
	}
	return nil
}

// Node is one entry in the path response: a level plus the user's state.
type Node struct {
	ID          string          `json:"id"`
	Kind        NodeKind        `json:"kind"`
	State       NodeState       `json:"state"`
	LessonCount int             `json:"lessonCount"`
	LessonsDone int             `json:"lessonsDone"`
	Payload     json.RawMessage `json:"payload"` // kind-discriminated
}

// Payloads of the kind-discriminated node union. Kinds without a payload
// type here (review, checkpoint, gate) carry an empty object.
type (
	// LessonPayload rides lesson and hard nodes.
	LessonPayload struct {
		XP int `json:"xp"`
	}
	// ChestPayload rides chest nodes; Opened renders the chest as claimed.
	ChestPayload struct {
		Gems   int  `json:"gems"`
		Opened bool `json:"opened,omitempty"`
	}
	// StoryPayload rides story nodes.
	StoryPayload struct {
		StoryID string `json:"storyId"`
	}
	// PracticePayload rides practice nodes; M6 fills the weak items.
	PracticePayload struct {
		ItemIDs []string `json:"itemIds"`
	}
)

// NodeProgress is the per-node user progress slice that state derivation
// reads: how many lessons are done, whether the node is complete (a chest
// counts as complete once opened, a jump gate completes skipped nodes), and
// whether the legendary tier is cleared.
type NodeProgress struct {
	LessonsDone int
	Completed   bool
	Legendary   bool
}

// DeriveStates resolves the icon state of every node on a path, in path
// order, from user progress. Completed and legendary come straight from
// progress; the first unfinished node is the single active one and every
// other unfinished node is locked. A fully finished path has no active node.
func DeriveStates(prog []NodeProgress) []NodeState {
	states := make([]NodeState, len(prog))
	activePlaced := false
	for i, p := range prog {
		switch {
		case p.Legendary:
			states[i] = StateLegendary
		case p.Completed:
			states[i] = StateCompleted
		case !activePlaced:
			states[i] = StateActive
			activePlaced = true
		default:
			states[i] = StateLocked
		}
	}
	return states
}
