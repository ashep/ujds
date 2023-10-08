package validation_test

import (
	"testing"

	"github.com/ashep/go-apperrors"
	"github.com/stretchr/testify/assert"

	"github.com/ashep/ujds/internal/validation"
)

func TestRecordIDValidator_Validate(tt *testing.T) {
	tt.Run("Empty", func(t *testing.T) {
		err := validation.NewRecordIDValidator().Validate("")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "record id",
			Reason: "must not be empty",
		})
	})
}
