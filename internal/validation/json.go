package validation

import (
	"encoding/json"
	"regexp"

	"github.com/ashep/go-apperrors"
	"github.com/xeipuuv/gojsonschema"
)

type JSONValidator struct {
	ldr map[*regexp.Regexp]gojsonschema.JSONLoader
}

func NewJSONValidator(schemas map[string]json.RawMessage) *JSONValidator {
	ldr := make(map[*regexp.Regexp]gojsonschema.JSONLoader)
	for pattern, sch := range schemas {
		re := regexp.MustCompile(pattern)
		ldr[re] = gojsonschema.NewBytesLoader(sch)
	}

	return &JSONValidator{
		ldr: ldr,
	}
}

func (v *JSONValidator) Validate(k, s string) error {
	if s == "" {
		return apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "empty",
		}
	}

	for re, ldr := range v.ldr {
		if !re.MatchString(k) {
			continue
		}

		res, err := gojsonschema.Validate(ldr, gojsonschema.NewBytesLoader([]byte(s)))
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
	}

	return nil
}
