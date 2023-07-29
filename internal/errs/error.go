package errs

import (
	"errors"
)

type NotFoundError struct {
	Subj string
}

func (e NotFoundError) Error() string {
	return e.Subj + " is not found"
}

func (e NotFoundError) Is(target error) bool {
	var notFoundError NotFoundError
	ok := errors.As(target, &notFoundError)

	return ok
}

type EmptyArgError struct {
	Subj string
}

func (e EmptyArgError) Error() string {
	return e.Subj + " is empty"
}

func (e EmptyArgError) Is(target error) bool {
	var emptyArgError EmptyArgError
	ok := errors.As(target, &emptyArgError)

	return ok
}

type InvalidArgError struct {
	Subj string
	E    error
}

func (e InvalidArgError) Error() string {
	s := ""

	if e.Subj != "" {
		s = "invalid " + e.Subj
		if e.E != nil {
			s += ": "
		}
	}

	if e.E != nil {
		s += e.E.Error()
	}

	return s
}

func (e InvalidArgError) Is(target error) bool {
	var invalidArgError InvalidArgError
	ok := errors.As(target, &invalidArgError)

	return ok
}

type AlreadyExistsError struct {
	Subj string
}

func (e AlreadyExistsError) Error() string {
	return e.Subj + " is already exists"
}

func (e AlreadyExistsError) Is(target error) bool {
	var alreadyExistsError AlreadyExistsError
	ok := errors.As(target, &alreadyExistsError)

	return ok
}

type AccessDeniedError struct{}

func (e AccessDeniedError) Error() string {
	return "access denied"
}

func (e AccessDeniedError) Is(target error) bool {
	var accessDeniedError AccessDeniedError
	ok := errors.As(target, &accessDeniedError)

	return ok
}
