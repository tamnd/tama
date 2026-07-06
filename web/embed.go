// Package web embeds the built frontend. Run `make web` (or `npm run build`
// in this directory) to refresh dist before building the binary.
package web

import "embed"

//go:embed all:dist
var Dist embed.FS
