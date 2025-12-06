package validation

import (
	"github.com/ashep/go-apperrors"
	"github.com/xeipuuv/gojsonschema"
)

type JSONValidator struct {
	ldr gojsonschema.JSONLoader
}

func NewJSONValidator(sch string) *JSONValidator {
	if len(sch) == 0 {
		sch = "{}"
	}

	return &JSONValidator{
		ldr: gojsonschema.NewBytesLoader([]byte(sch)),
	}
}

func (v *JSONValidator) Validate(data string) error {
	if v.ldr == nil {
		return nil
	}

	if data == "" {
		return apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "empty",
		}
	}

	res, err := gojsonschema.Validate(v.ldr, gojsonschema.NewBytesLoader([]byte(data)))
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
