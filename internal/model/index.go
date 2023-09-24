package model

import (
	"database/sql"
	"time"
)

type Index struct {
	ID        uint64
	Name      string
	Title     sql.NullString
	Schema    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
