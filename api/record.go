package api

import (
	"context"
	"fmt"

	"github.com/ashep/ujds/errs"
)

type Record struct {
	ID      string
	Schema  string
	Version uint64
	Data    string
}

func (a *API) PushRecords(ctx context.Context, items []Record) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	q := `INSERT INTO item (id, schema_id, data, version) VALUES ($1, $2, $3, $3)`
	for i, item := range items {
		if item.ID == "" {
			return errs.ErrEmptyArg{Subj: fmt.Sprintf("item id (%d)", i)}
		}
		if item.Schema == "" {
			return errs.ErrEmptyArg{Subj: fmt.Sprintf("item schmea (%d)", i)}
		}
		if item.Data == "" {
			return errs.ErrEmptyArg{Subj: fmt.Sprintf("item data (%d)", i)}
		}

		sch, err := a.GetSchema(ctx, item.Schema)
		if err != nil {
			return err
		}

		// TODO: validate data
		data := item.Data

		_, err = tx.ExecContext(ctx, q, item.ID, sch.ID, data)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (a *API) GetRecords(ctx context.Context, schema string, cursor uint64, limit uint32) ([]Record, uint64, error) {
	if limit == 0 || limit > 500 {
		limit = 500
	}

	sch, err := a.GetSchema(ctx, schema)
	if err != nil {
		return nil, 0, err
	}

	q := `SELECT id, version, data FROM item WHERE id>$1 AND schema_id=$2 AND version > $3 ORDER BY version LIMIT $4`
	rows, err := a.db.QueryContext(ctx, q, cursor, sch.ID, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	r := make([]Record, 0)
	id, ver, data := "", uint64(0), ""
	for rows.Next() {
		if err := rows.Scan(&id, &ver, &data); err != nil {
			return nil, 0, err
		}

		r = append(r, Record{
			ID:      id,
			Schema:  schema,
			Version: ver,
			Data:    data,
		})
	}

	return r, ver + 1, nil
}
