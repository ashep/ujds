package api_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/api"
	"github.com/ashep/ujds/internal/apperrors"
)

func TestIndex_Validate(tt *testing.T) {
	tt.Parallel()

	tt.Run("OkNilSchema", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, (&api.Index{}).Validate([]byte(`{"foo":"bar"}`)))
	})

	tt.Run("OkEmptySchema", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, (&api.Index{Schema: []byte("{}")}).Validate([]byte(`{"foo":"bar"}`)))
	})

	tt.Run("ErrInvalidSchema", func(t *testing.T) {
		t.Parallel()
		err := (&api.Index{Schema: []byte("{]")}).Validate([]byte(`{"foo":"bar"}`))
		require.EqualError(t, err, "schema validate failed: invalid character ']' looking for beginning of object key string")
	})

	tt.Run("ErrValidationError", func(t *testing.T) {
		t.Parallel()
		i := &api.Index{Schema: []byte(`{"properties":{"foo":{"type":"number"}}}`)}
		err := i.Validate([]byte(`{"foo":"bar"}`))
		require.EqualError(t, err, "foo: Invalid type. Expected: number, given: string")
	})

	tt.Run("ErrValidationError", func(t *testing.T) {
		t.Parallel()
		i := &api.Index{Schema: []byte(`{"properties":{"foo":{"type":"string"}}}`)}
		err := i.Validate([]byte(`{"foo":"bar"}`))
		require.NoError(t, err)
	})
}

func TestAPI_UpsertIndex(tt *testing.T) {
	tt.Parallel()

	tt.Run("ErrEmptyName", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		a := api.New(db, zerolog.Nop())
		err = a.UpsertIndex(context.Background(), "", "")

		require.ErrorIs(t, err, apperrors.EmptyArgError{Subj: "name"})
		assert.EqualError(t, err, "name is empty")
	})

	tt.Run("ErrInvalidSchema", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		require.NoError(t, err)

		a := api.New(db, zerolog.Nop())
		err = a.UpsertIndex(context.Background(), "theIndex", "{]")

		require.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "schema",
			Reason: "invalid character ']' looking for beginning of object key string",
		})
	})
}
