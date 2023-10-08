package validation_test

import (
	"testing"

	"github.com/ashep/go-apperrors"
	"github.com/stretchr/testify/assert"

	"github.com/ashep/ujds/internal/validation"
)

func TestIndexNameValidator_Validate(tt *testing.T) {
	tt.Run("Empty", func(t *testing.T) {
		err := validation.NewIndexNameValidator().Validate("")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not be empty",
		})
	})

	tt.Run("RegExpNotMatches", func(t *testing.T) {
		err := validation.NewIndexNameValidator().Validate("@bc")

		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must match the regexp ^[a-zA-Z0-9.-]{1,255}$",
		})
	})

	tt.Run("DotPrefix", func(t *testing.T) {
		err := validation.NewIndexNameValidator().Validate(".foo")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not start or end with a dot",
		})
	})

	tt.Run("DotSuffix", func(t *testing.T) {
		err := validation.NewIndexNameValidator().Validate("foo.")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not start or end with a dot",
		})
	})

	tt.Run("DashPrefix", func(t *testing.T) {
		err := validation.NewIndexNameValidator().Validate("-foo")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not start or end with a dash",
		})
	})

	tt.Run("DashSuffix", func(t *testing.T) {
		err := validation.NewIndexNameValidator().Validate("foo-")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not start or end with a dash",
		})
	})

	tt.Run("ConsecutiveDots", func(t *testing.T) {
		err := validation.NewIndexNameValidator().Validate("foo..bar")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not contain consecutive dots",
		})
	})
}
