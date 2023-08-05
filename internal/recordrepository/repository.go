package recordrepository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/indexrepository"
)

type Repository struct {
	db *sql.DB
	l  zerolog.Logger
}

func New(db *sql.DB, l zerolog.Logger) *Repository {
	return &Repository{db: db, l: l}
}

//nolint:cyclop // TODO
func (r *Repository) Push(ctx context.Context, index indexrepository.Index, records []Record) error {
	var err error

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("db begin failed: %w", err)
	}

	qGetRecord, err := tx.PrepareContext(ctx, `SELECT log_id FROM record WHERE checksum=$1`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("db prepare failed: %w", err)
	}

	defer func() {
		if err := qGetRecord.Close(); err != nil {
			r.l.Error().Err(err).Msg("prepared statement close failed")
		}
	}()

	qInsertLog, err := tx.PrepareContext(ctx, `INSERT INTO record_log (index_id, record_id, data)
		VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("db prepare failed: %w", err)
	}

	defer func() {
		if err := qInsertLog.Close(); err != nil {
			r.l.Error().Err(err).Msg("prepared statement close failed")
		}
	}()

	qInsertRecord, err := tx.PrepareContext(ctx, `INSERT INTO record (id, index_id, log_id, checksum)
VALUES ($1, $2, $3, $4) ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, updated_at=now()`)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("db prepare failed: %w", err)
	}

	defer func() {
		if err := qInsertRecord.Close(); err != nil {
			r.l.Error().Err(err).Msg("prepared statement close failed")
		}
	}()

	for i, rec := range records {
		if rec.ID == "" {
			_ = tx.Rollback()
			return apperrors.EmptyArgError{Subj: fmt.Sprintf("record %d: id", i)}
		}

		if rec.Data == "" {
			_ = tx.Rollback()
			return apperrors.EmptyArgError{Subj: fmt.Sprintf("record %d: data", i)}
		}

		// Validate data against schema
		recDataB := []byte(rec.Data)
		if err = index.Validate(recDataB); err != nil {
			_ = tx.Rollback()
			return apperrors.InvalidArgError{Subj: fmt.Sprintf("record data (%d)", i), Reason: err.Error()}
		}

		logID := uint64(0)

		sumSrc := append(recDataB, []byte(rec.Index)...) //nolint:gocritic // it's ok
		sumSrc = append(sumSrc, []byte(rec.ID)...)
		sum := sha256.Sum256(sumSrc)

		// Check if we already have such data recorded as latest version
		row := qGetRecord.QueryRowContext(ctx, sum[:])
		if err = row.Scan(&logID); errors.Is(err, sql.ErrNoRows) { //nolint:revive // this is intentional empty block
			// Ok, continue to insert
		} else if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("db scan failed: %w", err)
		} else {
			// A record with the same data found, skip it
			continue
		}

		row = qInsertLog.QueryRowContext(ctx, index.ID, rec.ID, rec.Data)
		if err = row.Scan(&logID); err != nil && errors.Is(err, sql.ErrNoRows) {
			_ = tx.Rollback()
			return nil
		} else if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("db query failed: %w", err)
		}

		if _, err = qInsertRecord.ExecContext(ctx, rec.ID, index.ID, logID, sum[:]); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("db query failed: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db commit failed: %w", err)
	}

	return nil
}

// Get returns last version of a record.
func (r *Repository) Get(ctx context.Context, indexName string, id string) (Record, error) {
	q := `SELECT r.log_id, l.data, r.created_at, r.updated_at FROM record r
		LEFT JOIN record_log l ON r.log_id = l.id
		LEFT JOIN index i ON r.index_id = i.id
		WHERE i.name=$1 AND r.id=$2 ORDER BY l.created_at DESC LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, indexName, id)

	rec := Record{
		ID:    id,
		Index: indexName,
	}

	err := row.Scan(&rec.Rev, &rec.Data, &rec.CreatedAt, &rec.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Record{}, apperrors.NotFoundError{Subj: "record"}
	} else if err != nil {
		return Record{}, fmt.Errorf("db scan failed: %w", err)
	}

	return rec, nil
}

func (r *Repository) GetAll(ctx context.Context, indexName string, since time.Time, cursor uint64, limit uint32) ([]Record, uint64, error) {
	if limit == 0 || limit > 500 {
		limit = 500
	}

	q := `SELECT r.id, r.log_id, l.data, r.created_at, r.updated_at FROM record r
		LEFT JOIN record_log l ON r.log_id = l.id
		LEFT JOIN index i ON r.index_id = i.id
		WHERE i.name=$1 AND r.updated_at >= $2 AND l.id >= $3 ORDER BY l.id LIMIT $4`

	rows, err := r.db.QueryContext(ctx, q, indexName, since, cursor, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("db query failed: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	rcs := make([]Record, 0)
	recID, logID, data, crAt, upAt := "", uint64(0), "", time.Time{}, time.Time{}

	for rows.Next() {
		if err := rows.Scan(&recID, &logID, &data, &crAt, &upAt); err != nil {
			return nil, 0, fmt.Errorf("db scan failed: %w", err)
		}

		rcs = append(rcs, Record{
			ID:        recID,
			Index:     indexName,
			Rev:       logID,
			Data:      data,
			CreatedAt: crAt,
			UpdatedAt: upAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("db rows iteration failed: %w", err)
	}

	nextCursor := uint64(0)
	if len(rcs) > 0 {
		nextCursor = logID + 1
	}

	return rcs, nextCursor, nil
}

func (r *Repository) Clear(ctx context.Context, indexName string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM record WHERE index_id=$1`, indexName)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete records: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM record_log WHERE index_id=$1`, indexName)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete records: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db commit failed: %w", err)
	}

	return nil
}