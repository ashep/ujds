package validation

import (
	"github.com/ashep/go-apperrors"
	"github.com/xeipuuv/gojsonschema"
)

type JSONValidator struct{}

func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

func (v *JSONValidator) Validate(schema, data []byte) error {
	if len(data) == 0 {
		return apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "empty",
		}
	}

	if schema == nil {
		schema = []byte("{}")
	}

	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(schema), gojsonschema.NewBytesLoader(data))
	if err != nil {
		return apperrors.InvalidArgError{
			Subj:   "json schema or data",
			Reason: err.Error(),
		}
	}

	if !res.Valid() {
		return apperrors.InvalidArgError{
			Subj:   "json",
			Reason: res.Errors()[0].String(),
		}
	}

	return nil
}
