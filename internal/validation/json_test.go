package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/validation"
)

func Test_ValidateJSON(tt *testing.T) {
	tt.Run("EmptyJSON", func(t *testing.T) {
		assert.EqualError(t, validation.NewJSONValidator().Validate(nil, nil), "invalid json: empty")
	})

	tt.Run("NilSchema", func(t *testing.T) {
		require.NoError(t, validation.NewJSONValidator().Validate(nil, []byte(`{"foo":"bar"}`)))
	})

	tt.Run("MalformedSchema", func(t *testing.T) {
		err := validation.NewJSONValidator().Validate([]byte("{]"), []byte(`{"foo":"bar"}`))
		require.EqualError(t, err, "invalid json schema or data: invalid character ']' looking for beginning of object key string")
	})

	tt.Run("MalformedData", func(t *testing.T) {
		err := validation.NewJSONValidator().Validate(nil, []byte(`{]`))
		require.EqualError(t, err, "invalid json schema or data: invalid character ']' looking for beginning of object key string")
	})

	tt.Run("DataValidationError", func(t *testing.T) {
		err := validation.NewJSONValidator().Validate([]byte(`{"properties":{"foo":{"type":"number"}}}`), []byte(`{"foo":"bar"}`))
		require.EqualError(t, err, "invalid json: foo: Invalid type. Expected: number, given: string")
	})

	tt.Run("Ok", func(t *testing.T) {
		err := validation.NewJSONValidator().Validate([]byte(`{"properties":{"foo":{"type":"string"}}}`), []byte(`{"foo":"bar"}`))
		require.NoError(t, err)
	})
}
