package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/xeipuuv/gojsonschema"

	"github.com/ashep/ujds/errs"
)

type Schema struct {
	Id        int
	Name      string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *API) UpsertSchema(ctx context.Context, name, data string) error {
	if name == "" {
		return errs.ErrEmptyArg{Subj: "name"}
	}
	if data == "" {
		return errs.ErrEmptyArg{Subj: "data"}
	}

	if err := json.Unmarshal([]byte(data), &struct{}{}); err != nil {
		return errs.ErrInvalidArg{Subj: "data", E: err}
	}

	q := `INSERT INTO schema (name, data) VALUES ($1, $2) 
ON CONFLICT (name) DO UPDATE SET data=$2, updated_at=now()`
	if _, err := a.db.ExecContext(ctx, q, name, data); err != nil {
		return err
	}

	return nil
}

func (a *API) GetSchema(ctx context.Context, name string) (*Schema, error) {
	var (
		id        int
		data      []byte
		createdAt time.Time
		updatedAt time.Time
	)

	q := `SELECT id, data, created_at, updated_at FROM schema WHERE name=$1`
	row := a.db.QueryRowContext(ctx, q, name)
	if err := row.Scan(&id, &data, &createdAt, &updatedAt); errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound{Subj: "schema"}
	} else if err != nil {
		return nil, err
	}

	return &Schema{Id: id, Name: name, Data: data, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}

func (s *Schema) Validate(data []byte) error {
	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(s.Data), gojsonschema.NewBytesLoader(data))
	if err != nil {
		return err
	}

	if !res.Valid() {
		return errors.New(res.Errors()[0].String())
	}

	return nil
}
