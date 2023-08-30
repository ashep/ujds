package model

import (
	"time"
)

type RecordUpdate struct {
	ID      string
	IndexID uint64
	Data    string
}

type Record struct {
	ID        string
	IndexID   uint64
	Rev       uint64
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
