package validation_test

import (
	"testing"

	"github.com/ashep/go-apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/validation"
)

func Test_ValidateJSON(tt *testing.T) {
	tt.Run("EmptySchema", func(t *testing.T) {
		v := validation.NewJSONValidator("")
		require.NotNil(t, v)
		// Empty schema ("") should use default schema
		err := v.Validate(`{"foo":"bar"}`)
		assert.NoError(t, err)
	})

	tt.Run("EmptyJSON", func(t *testing.T) {
		v := validation.NewJSONValidator(`{}`)

		err := v.Validate("")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "empty",
		})
	})

	tt.Run("MalformedSchema", func(t *testing.T) {
		v := validation.NewJSONValidator("{]")

		err := v.Validate(`{"foo":"bar"}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json schema or data",
			Reason: "invalid character ']' looking for beginning of object key string",
		})
	})

	tt.Run("MalformedData", func(t *testing.T) {
		v := validation.NewJSONValidator(`{}`)

		err := v.Validate(`{]`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json schema or data",
			Reason: "invalid character ']' looking for beginning of object key string",
		})
	})

	tt.Run("DataValidationError", func(t *testing.T) {
		v := validation.NewJSONValidator(`{"properties":{"foo":{"type":"number"}}}`)

		err := v.Validate(`{"foo":"bar"}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "foo: Invalid type. Expected: number, given: string",
		})
	})

	tt.Run("Ok", func(t *testing.T) {
		v := validation.NewJSONValidator(`{"properties":{"foo":{"type":"string"}}}`)

		err := v.Validate(`{"foo":"bar"}`)
		assert.NoError(t, err)
	})
}
