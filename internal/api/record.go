package api

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/ujds/internal/errs"
)

type Record struct {
	Id        string
	Index     string
	Rev       uint64
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *API) PushRecords(ctx context.Context, schema string, records []Record) error {
	var err error

	sch, err := a.GetIndex(ctx, schema)
	if err != nil {
		return err
	}

	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	qGetRecord, err := tx.PrepareContext(ctx, `SELECT log_id FROM record WHERE checksum=$1`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	qInsertLog, err := tx.PrepareContext(ctx, `INSERT INTO record_log (index_id, record_id, data) 
		VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	qInsertRecord, err := tx.PrepareContext(ctx, `INSERT INTO record (id, index_id, log_id, checksum) 
VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for i, rec := range records {
		if rec.Id == "" {
			_ = tx.Rollback()
			return errs.ErrEmptyArg{Subj: fmt.Sprintf("record %d: id", i)}
		}
		if rec.Data == "" {
			_ = tx.Rollback()
			return errs.ErrEmptyArg{Subj: fmt.Sprintf("record %d: data", i)}
		}

		// Validate data against schema
		recDataB := []byte(rec.Data)
		if err = sch.Validate(recDataB); err != nil {
			_ = tx.Rollback()
			return errs.ErrInvalidArg{Subj: fmt.Sprintf("record data (%d)", i), E: err}
		}

		// Check if we already have such data recorded as latest version
		logId := uint64(0)
		sumSrc := append(recDataB, []byte(rec.Index)...)
		sumSrc = append(sumSrc, []byte(rec.Id)...)
		sum := sha256.Sum256(sumSrc)
		row := qGetRecord.QueryRowContext(ctx, sum[:])
		if err = row.Scan(&logId); errors.Is(err, sql.ErrNoRows) {
			// Ok, continue to insert
		} else if err != nil {
			_ = tx.Rollback()
			return err
		} else {
			// A record with the same data found, skip it
			continue
		}

		row = qInsertLog.QueryRowContext(ctx, sch.Id, rec.Id, rec.Data)
		if err = row.Scan(&logId); err != nil && errors.Is(err, sql.ErrNoRows) {
			_ = tx.Rollback()
			return nil
		} else if err != nil {
			_ = tx.Rollback()
			return err
		}

		if _, err = qInsertRecord.ExecContext(ctx, rec.Id, sch.Id, logId, sum[:]); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (a *API) GetRecords(
	ctx context.Context,
	schema string,
	since time.Time,
	cursor uint64,
	limit uint32,
) ([]Record, uint64, error) {
	if limit == 0 || limit > 500 {
		limit = 500
	}

	q := `SELECT r.id, r.log_id, l.data, r.created_at, r.updated_at FROM record r
		LEFT JOIN record_log l ON r.log_id = l.id 
		LEFT JOIN index s ON r.index_id = s.id 
		WHERE s.name=$1 AND r.updated_at >= $2 AND l.id >= $3 ORDER BY l.id LIMIT $4`
	rows, err := a.db.QueryContext(ctx, q, schema, since, cursor, limit)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = rows.Close()
	}()

	r := make([]Record, 0)
	recId, logId, data, crAt, upAt := "", uint64(0), "", time.Time{}, time.Time{}
	for rows.Next() {
		if err := rows.Scan(&recId, &logId, &data, &crAt, &upAt); err != nil {
			return nil, 0, err
		}

		r = append(r, Record{
			Id:        recId,
			Index:     schema,
			Rev:       logId,
			Data:      data,
			CreatedAt: crAt,
			UpdatedAt: upAt,
		})
	}

	nextCursor := uint64(0)
	if len(r) > 0 {
		nextCursor = logId + 1
	}

	return r, nextCursor, nil
}
