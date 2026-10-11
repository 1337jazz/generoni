// Package web holds the HTML templates and static files, embedded into the
// binary so the image needs no files on disk.
package web

import "embed"

//go:embed templates
var Templates embed.FS

// Paths keep the "static/" prefix, so a request for /static/style.css maps
// straight to static/style.css with no StripPrefix needed.
//
//go:embed static
var Static embed.FS
