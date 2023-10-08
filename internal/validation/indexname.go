package validation

import (
	"regexp"
	"strings"

	"github.com/ashep/go-apperrors"
)

type IndexNameValidator struct {
	nameRe *regexp.Regexp
}

func NewIndexNameValidator() *IndexNameValidator {
	return &IndexNameValidator{
		nameRe: regexp.MustCompile("^[a-zA-Z0-9.-]{1,255}$"),
	}
}

func (v *IndexNameValidator) Validate(s string) error {
	if s == "" {
		return apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not be empty",
		}
	}

	if !v.nameRe.MatchString(s) {
		return apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must match the regexp " + v.nameRe.String(),
		}
	}

	if strings.HasPrefix(s, ".") || strings.HasSuffix(s, ".") {
		return apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not start or end with a dot",
		}
	}

	if strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") {
		return apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not start or end with a dash",
		}
	}

	if strings.Contains(s, "..") {
		return apperrors.InvalidArgError{
			Subj:   "index name",
			Reason: "must not contain consecutive dots",
		}
	}

	return nil
}
