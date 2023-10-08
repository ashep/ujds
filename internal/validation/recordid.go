package validation

import (
	"github.com/ashep/go-apperrors"
)

type RecordIDValidator struct{}

func NewRecordIDValidator() *RecordIDValidator {
	return &RecordIDValidator{}
}

func (v *RecordIDValidator) Validate(s string) error {
	if s == "" {
		return apperrors.InvalidArgError{
			Subj:   "record id",
			Reason: "must not be empty",
		}
	}

	return nil
}
