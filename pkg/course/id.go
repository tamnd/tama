package course

import (
	"fmt"
	"strings"

	"github.com/tamnd/tama/pkg/course/lang"
)

// idSep splits a course ID into target and base.
const idSep = "-from-"

// InvalidIDError is the typed rejection for a malformed course ID.
type InvalidIDError struct {
	ID     string
	Reason string
}

func (e *InvalidIDError) Error() string {
	return fmt.Sprintf("course: invalid id %q: %s", e.ID, e.Reason)
}

// SamePairError is the typed rejection for IDs like "en-from-en" where the
// target and base resolve to the same language.
type SamePairError struct {
	ID   string
	Lang lang.Language
}

func (e *SamePairError) Error() string {
	return fmt.Sprintf("course: id %q pairs %s with itself", e.ID, e.Lang.Name)
}

// ID builds the canonical course ID, "<target>-from-<base>" using each
// language's shortest code: ja-from-en, en-from-zh, qhv-from-en.
func ID(target, base lang.Language) string {
	return target.Code() + idSep + base.Code()
}

// ParseID resolves a course ID into its target and base languages. Both
// sides go through lang.Lookup, so 639-3 codes and aliases work too, and the
// canonical short form comes back from ID. Errors are typed: *InvalidIDError
// for a malformed ID, *lang.NotFoundError for an unknown code, and
// *SamePairError when both sides are the same language.
func ParseID(id string) (target, base lang.Language, err error) {
	targetCode, baseCode, ok := strings.Cut(id, idSep)
	if !ok || targetCode == "" || baseCode == "" {
		return target, base, &InvalidIDError{ID: id, Reason: "want <target>-from-<base>"}
	}
	if target, err = lang.Lookup(targetCode); err != nil {
		return target, base, err
	}
	if base, err = lang.Lookup(baseCode); err != nil {
		return target, base, err
	}
	if target.Code639_3 == base.Code639_3 {
		return target, base, &SamePairError{ID: id, Lang: target}
	}
	return target, base, nil
}
