package migration

import "embed"

// FS embeds the SQL migration files so they are compiled into the binary and
// can be applied at startup without shipping the folder or the migrate CLI.
//
//go:embed *.sql
var FS embed.FS
