package errs

type ErrNotFound struct {
	Subj string
}

func (e ErrNotFound) Error() string {
	return e.Subj + " is not found"
}

func (e ErrNotFound) Is(target error) bool {
	_, ok := target.(ErrNotFound)
	return ok
}

type ErrEmptyArg struct {
	Subj string
}

func (e ErrEmptyArg) Error() string {
	return e.Subj + " is empty"
}

func (e ErrEmptyArg) Is(target error) bool {
	_, ok := target.(ErrEmptyArg)
	return ok
}

type ErrInvalidArg struct {
	Subj string
	E    error
}

func (e ErrInvalidArg) Error() string {
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

func (e ErrInvalidArg) Is(target error) bool {
	_, ok := target.(ErrInvalidArg)
	return ok
}

type ErrAlreadyExists struct {
	Subj string
}

func (e ErrAlreadyExists) Error() string {
	return e.Subj + " is already exists"
}

func (e ErrAlreadyExists) Is(target error) bool {
	_, ok := target.(ErrAlreadyExists)
	return ok
}

type ErrAccessDenied struct {
}

func (e ErrAccessDenied) Error() string {
	return "access denied"
}

func (e ErrAccessDenied) Is(target error) bool {
	_, ok := target.(ErrAccessDenied)
	return ok
}
