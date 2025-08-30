package indexrepo

import (
	"database/sql"
	"time"
)

type IndexFilter struct {
	Names []string
}

type Index struct {
	ID        uint64
	Name      string
	Title     sql.NullString
	Schema    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
