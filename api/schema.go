package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/xeipuuv/gojsonschema"

	"github.com/ashep/ujds/errs"
)

type Schema struct {
	ID        int
	Name      string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *API) PushSchema(ctx context.Context, name, schema string) error {
	if name == "" {
		return errs.ErrEmptyArg{Subj: "name"}
	}
	if schema == "" {
		return errs.ErrEmptyArg{Subj: "schema"}
	}

	if err := json.Unmarshal([]byte(schema), &struct{}{}); err != nil {
		return errs.ErrInvalidArg{Subj: "json data", E: err}
	}

	q := `INSERT INTO schema (name, version, data) 
VALUES ($1, nextval('item_schema_version_seq'), $2) RETURNING version`
	row := a.db.QueryRowContext(ctx, q, name, schema)

	v := uint32(0)
	if err := row.Scan(&v); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return errs.ErrAlreadyExists{Subj: "schema"}
		}
		return err
	}

	return nil
}

func (a *API) GetSchema(ctx context.Context, name string) (Schema, error) {
	var (
		id     int
		schema []byte
	)

	q := `SELECT id, data FROM schema WHERE name=$1 LIMIT 1`
	row := a.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&id, &schema); errors.Is(err, sql.ErrNoRows) {
		return Schema{}, errs.ErrNotFound{Subj: "schema"}
	} else if err != nil {
		return Schema{}, nil
	}

	return Schema{ID: id, Name: name, Data: schema}, nil
}

func (a *API) UpdateSchema(ctx context.Context, name, data string, ver uint32) (uint32, error) {
	if ver == 0 {
		return 0, errs.ErrEmptyArg{Subj: "version"}
	}

	if err := json.Unmarshal([]byte(data), &struct{}{}); err != nil {
		return ver, errs.ErrInvalidArg{Subj: "json data", E: err}
	}

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return ver, err
	}

	q := `SELECT version FROM schema WHERE name=$1 AND version=$2 FOR UPDATE`
	row := tx.QueryRowContext(ctx, q, name, ver)
	if err := row.Scan(&ver); errors.Is(err, sql.ErrNoRows) {
		_ = tx.Rollback()
		return ver, errs.ErrNotFound{Subj: "schema"}
	} else if err != nil {
		_ = tx.Rollback()
		return ver, err
	}

	q = `UPDATE schema SET data=$1, version=nextval('item_schema_version_seq') 
WHERE name=$2 AND version=$3 RETURNING version`
	row = tx.QueryRowContext(ctx, q, data, name, ver)
	if err := row.Scan(&ver); err != nil {
		_ = tx.Rollback()
		return ver, err
	}

	if err := tx.Commit(); err != nil {
		return ver, err
	}

	a.mux.Lock()
	delete(a.schemaCache, name)
	a.mux.Unlock()

	return ver, nil
}

func (s Schema) Validate(data []byte) error {
	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(s.Data), gojsonschema.NewBytesLoader(data))
	if err != nil {
		return err
	}

	if !res.Valid() {
		return errors.New(res.Errors()[0].String())
	}

	return nil
}
