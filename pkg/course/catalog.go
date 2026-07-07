// Package course defines the course model: the language registry under
// lang/, course IDs and the seed catalog, the path hierarchy (sections,
// units, levels, lessons) with its node kinds and states, guidebooks,
// stories, and the frozen pack format the generator produces. Any base to
// any target is a valid course; the seed catalog is just what the picker
// surfaces first.
package course

import "github.com/tamnd/tama/pkg/course/lang"

// Languages is the full registry sorted by English name, the shape the
// languages endpoint serves. Every entry is usable as base or target.
var Languages = lang.All()
