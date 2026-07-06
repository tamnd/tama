package store

import _ "embed"

// demoPack is the handcrafted es-en pack used by `tama seed --demo`. It is a
// checked-in fixture, never model output, so the dev loop works offline.
//
//go:embed testdata/demo-pack.json
var demoPack []byte

// DemoPack returns the raw (uncompressed) demo pack JSON.
func DemoPack() []byte {
	return demoPack
}
