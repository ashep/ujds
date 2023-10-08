package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/model"
)

func Test_ValidateJSON(tt *testing.T) {
	tt.Run("NilSchema", func(t *testing.T) {
		require.NoError(t, model.ValidateJSON(nil, []byte(`{"foo":"bar"}`)))
	})

	tt.Run("EmptySchema", func(t *testing.T) {
		require.NoError(t, model.ValidateJSON([]byte(`{}`), []byte(`{"foo":"bar"}`)))
	})

	tt.Run("InvalidSchema", func(t *testing.T) {
		err := model.ValidateJSON([]byte("{]"), []byte(`{"foo":"bar"}`))
		require.EqualError(t, err, "schema validate failed: invalid character ']' looking for beginning of object key string")
	})

	tt.Run("ValidationError", func(t *testing.T) {
		err := model.ValidateJSON([]byte(`{"properties":{"foo":{"type":"number"}}}`), []byte(`{"foo":"bar"}`))
		require.EqualError(t, err, "foo: Invalid type. Expected: number, given: string")
	})

	tt.Run("Ok", func(t *testing.T) {
		err := model.ValidateJSON([]byte(`{"properties":{"foo":{"type":"string"}}}`), []byte(`{"foo":"bar"}`))
		require.NoError(t, err)
	})
}
