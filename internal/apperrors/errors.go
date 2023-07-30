package apperrors

type NotFoundError struct {
	Subj string
}

func (e NotFoundError) Error() string {
	return e.Subj + " is not found"
}

type AlreadyExistsError struct {
	Subj string
}

func (e AlreadyExistsError) Error() string {
	return e.Subj + " is already exists"
}

type AccessDeniedError struct{}

func (e AccessDeniedError) Error() string {
	return "access denied"
}

type EmptyArgError struct {
	Subj string
}

func (e EmptyArgError) Error() string {
	return e.Subj + " is empty"
}

type InvalidArgError struct {
	Subj   string
	Reason string
}

func (e InvalidArgError) Error() string {
	s := ""

	if e.Subj != "" {
		s = "invalid " + e.Subj
		if e.Reason != "" {
			s += ": "
		}
	}

	s += e.Reason

	return s
}
