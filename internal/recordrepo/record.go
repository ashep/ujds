package recordrepo

import (
	"crypto/sha256"
	"encoding/binary"
	"time"
)

type RecordUpdate struct {
	ID      string
	IndexID uint64
	Data    string
}

func (rec *RecordUpdate) Checksum() []byte {
	indexIDBytes := make([]byte, 8) //nolint:mnd // it's ok
	binary.LittleEndian.PutUint64(indexIDBytes, rec.IndexID)

	sumSrc := append([]byte(rec.Data), indexIDBytes...)
	sumSrc = append(sumSrc, []byte(rec.ID)...)
	sum := sha256.Sum256(sumSrc)

	return sum[:]
}

type Record struct {
	ID        string
	IndexID   uint64
	Rev       uint64
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
	TouchedAt time.Time
}
