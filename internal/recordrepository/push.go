package recordrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"

	"github.com/ashep/ujds/internal/model"
)

type statements struct {
	getLog       *sql.Stmt
	insertLog    *sql.Stmt
	upsertRecord *sql.Stmt
	touchRecord  *sql.Stmt
}

func (s *statements) Close(l zerolog.Logger) {
	if err := s.getLog.Close(); err != nil {
		l.Error().Err(err).Msg("prepared statement close failed")
	}

	if err := s.insertLog.Close(); err != nil {
		l.Error().Err(err).Msg("prepared statement close failed")
	}

	if err := s.upsertRecord.Close(); err != nil {
		l.Error().Err(err).Msg("prepared statement close failed")
	}

	if err := s.touchRecord.Close(); err != nil {
		l.Error().Err(err).Msg("prepared statement close failed")
	}
}

func (r *Repository) Push(ctx context.Context, updates []model.RecordUpdate) error {
	var err error

	if len(updates) == 0 {
		return apperrors.InvalidArgError{Subj: "updates", Reason: "must not be empty"}
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("db begin: %w", err)
	}

	stmt, err := r.prepareStatements(ctx, tx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("prepare statements: %w", err)
	}

	defer stmt.Close(r.l)

	for i, rec := range updates {
		if rec.IndexID == 0 {
			return apperrors.InvalidArgError{Subj: fmt.Sprintf("record %d", i), Reason: "zero index id"}
		}

		if err = r.upsertOrTouch(ctx, stmt, rec); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (r *Repository) prepareStatements(ctx context.Context, tx *sql.Tx) (*statements, error) {
	getLog, err := tx.PrepareContext(ctx, `SELECT log_id FROM record WHERE checksum=$1`)
	if err != nil {
		return nil, fmt.Errorf("get record by log id: %w", err)
	}

	insertLog, err := tx.PrepareContext(ctx, `INSERT INTO record_log (index_id, record_id, data)
		VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("insert record log: %w", err)
	}

	upsertRecord, err := tx.PrepareContext(ctx, `INSERT INTO record (id, index_id, log_id, checksum, data)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id, index_id) DO UPDATE SET log_id=$3, checksum=$4, data=$5, updated_at=now(), touched_at=now()`)
	if err != nil {
		return nil, fmt.Errorf("insert record: %w", err)
	}

	touchRecord, err := tx.PrepareContext(ctx, `UPDATE record SET touched_at=now() WHERE log_id=$1`)
	if err != nil {
		return nil, fmt.Errorf("update record touch time: %w", err)
	}

	return &statements{
		getLog:       getLog,
		insertLog:    insertLog,
		upsertRecord: upsertRecord,
		touchRecord:  touchRecord,
	}, nil
}

func (r *Repository) upsertOrTouch(ctx context.Context, stmt *statements, upd model.RecordUpdate) error {
	if err := r.recordIDValidator.Validate(upd.ID); err != nil {
		return err //nolint:wrapcheck // ok
	}

	if upd.Data == "" {
		return apperrors.InvalidArgError{Subj: "record data", Reason: "must not be empty"}
	}

	if err := r.jsonValidator.Validate(upd.Schema, []byte(upd.Data)); err != nil {
		return apperrors.InvalidArgError{Subj: "record data", Reason: err.Error()}
	}

	row := stmt.getLog.QueryRowContext(ctx, upd.Checksum())

	logID := uint64(0)
	if err := row.Scan(&logID); errors.Is(err, sql.ErrNoRows) {
		// There is no record with requested checksum exists, try to insert it
		if err := r.upsert(ctx, stmt.insertLog, stmt.upsertRecord, upd); err != nil {
			return fmt.Errorf("upsert record: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("get record by checksum scan: %w", err)
	} else {
		// There is a record with the same checksum exists, just touch it
		if err := r.touch(ctx, stmt.touchRecord, logID); err != nil {
			return fmt.Errorf("touch record: %w", err)
		}
	}

	return nil
}

// upsert
func (r *Repository) upsert(
	ctx context.Context,
	insertLogStmt *sql.Stmt,
	upsertRecordStmt *sql.Stmt,
	upd model.RecordUpdate,
) error {
	var logID uint64

	row := insertLogStmt.QueryRowContext(ctx, upd.IndexID, upd.ID, upd.Data)
	if err := row.Scan(&logID); err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		return fmt.Errorf("insert log db query: %w", err)
	}

	if _, err := upsertRecordStmt.ExecContext(ctx, upd.ID, upd.IndexID, logID, upd.Checksum(), upd.Data); err != nil {
		return fmt.Errorf("insert record db query: %w", err)
	}

	return nil
}

func (r *Repository) touch(ctx context.Context, stmt *sql.Stmt, logID uint64) error {
	res, err := stmt.ExecContext(ctx, logID)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get db rows affected: %w", err)
	}

	if ra == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
