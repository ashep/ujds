package api

import (
	"bytes"
	"errors"
	"time"

	"github.com/xeipuuv/gojsonschema"
)

type Schema struct {
	Id        int
	Name      string
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Schema) Validate(data []byte) error {
	if bytes.Equal(s.Data, []byte("{}")) {
		return nil
	}

	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(s.Data), gojsonschema.NewBytesLoader(data))
	if err != nil {
		return err
	}

	if !res.Valid() {
		return errors.New(res.Errors()[0].String())
	}

	return nil
}
