package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/model"
)

func TestIndex_Validate(tt *testing.T) {
	tt.Parallel()

	tt.Run("NilSchema", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, (&model.Index{}).Validate([]byte(`{"foo":"bar"}`)))
	})

	tt.Run("EmptySchema", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, (&model.Index{Schema: []byte("{}")}).Validate([]byte(`{"foo":"bar"}`)))
	})

	tt.Run("InvalidSchema", func(t *testing.T) {
		t.Parallel()
		err := (&model.Index{Schema: []byte("{]")}).Validate([]byte(`{"foo":"bar"}`))
		require.EqualError(t, err, "schema validate failed: invalid character ']' looking for beginning of object key string")
	})

	tt.Run("ValidationError", func(t *testing.T) {
		t.Parallel()
		i := &model.Index{Schema: []byte(`{"properties":{"foo":{"type":"number"}}}`)}
		err := i.Validate([]byte(`{"foo":"bar"}`))
		require.EqualError(t, err, "foo: Invalid type. Expected: number, given: string")
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()
		i := &model.Index{Schema: []byte(`{"properties":{"foo":{"type":"string"}}}`)}
		err := i.Validate([]byte(`{"foo":"bar"}`))
		require.NoError(t, err)
	})
}
