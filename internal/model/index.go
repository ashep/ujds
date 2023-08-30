package model

import (
	"time"
)

type Index struct {
	ID        uint64
	Name      string
	Schema    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
