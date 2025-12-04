package sql

import (
	"embed"
)

//go:embed migrations/*.sql
var FS embed.FS
