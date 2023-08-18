package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/xeipuuv/gojsonschema"
)

type Index struct {
	ID        uint
	Name      string
	Schema    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Index) Validate(data []byte) error {
	if len(data) != 0 && !json.Valid(data) {
		return errors.New("invalid json")
	}

	if len(s.Schema) == 0 || bytes.Equal(s.Schema, []byte("{}")) {
		return nil
	}

	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(s.Schema), gojsonschema.NewBytesLoader(data))
	if err != nil {
		return fmt.Errorf("schema validate failed: %w", err)
	}

	if !res.Valid() {
		return errors.New(res.Errors()[0].String())
	}

	return nil
}
