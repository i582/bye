package static

import "embed"

//go:embed styles templates config hello-world pages index.md
var Content embed.FS
