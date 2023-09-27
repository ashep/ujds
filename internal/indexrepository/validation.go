package indexrepository

import (
	"errors"
	"regexp"
	"strings"
)

type NameValidator struct {
	nameRe *regexp.Regexp
}

func NewNameValidator() *NameValidator {
	return &NameValidator{
		nameRe: regexp.MustCompile("^[a-zA-Z0-9.-]{1,255}$"),
	}
}

func (v *NameValidator) Validate(s string) error {
	if s == "" {
		return errors.New("must not be empty")
	}

	if !v.nameRe.MatchString(s) {
		return errors.New("must match the regexp " + v.nameRe.String())
	}

	if strings.HasPrefix(s, ".") || strings.HasSuffix(s, ".") {
		return errors.New("must not start or end with a dot")
	}

	if strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") {
		return errors.New("must not start or end with a dash")
	}

	if strings.Contains(s, "..") {
		return errors.New("must not contain consecutive dots")
	}

	return nil
}
