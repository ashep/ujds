package recordrepository

import (
	"time"
)

type Record struct {
	ID        string
	Index     string
	Rev       uint64
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
