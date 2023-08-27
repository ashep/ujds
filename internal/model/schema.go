package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateJSON(schema, data []byte) error {
	if len(data) != 0 && !json.Valid(data) {
		return errors.New("invalid json")
	}

	if len(schema) == 0 || bytes.Equal(schema, []byte("{}")) {
		return nil
	}

	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(schema), gojsonschema.NewBytesLoader(data))
	if err != nil {
		return fmt.Errorf("schema validate failed: %w", err)
	}

	if !res.Valid() {
		return errors.New(res.Errors()[0].String())
	}

	return nil
}
