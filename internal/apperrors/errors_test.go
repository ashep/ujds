package apperrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ashep/ujds/internal/apperrors"
)

func TestNotFoundError(t *testing.T) {
	err := apperrors.NotFoundError{Subj: "foo"}

	assert.Equal(t, "foo is not found", err.Error())

	assert.True(t, errors.Is(err, apperrors.NotFoundError{Subj: "foo"}))
	assert.True(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.NotFoundError{Subj: "foo"}))

	assert.False(t, errors.Is(err, apperrors.NotFoundError{Subj: "bar"}))
	assert.False(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.NotFoundError{Subj: "bar"}))

	assert.True(t, errors.As(err, &apperrors.NotFoundError{}))
	assert.True(t, errors.As(fmt.Errorf("wrap: %w", err), &apperrors.NotFoundError{}))
}

func TestAlreadyExistsError(t *testing.T) {
	err := apperrors.AlreadyExistsError{Subj: "foo"}

	assert.Equal(t, "foo is already exists", err.Error())

	assert.True(t, errors.Is(err, apperrors.AlreadyExistsError{Subj: "foo"}))
	assert.True(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.AlreadyExistsError{Subj: "foo"}))

	assert.False(t, errors.Is(err, apperrors.AlreadyExistsError{Subj: "bar"}))
	assert.False(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.AlreadyExistsError{Subj: "bar"}))

	assert.True(t, errors.As(err, &apperrors.AlreadyExistsError{}))
	assert.True(t, errors.As(fmt.Errorf("wrap: %w", err), &apperrors.AlreadyExistsError{}))
}

func TestEmptyArgError(t *testing.T) {
	err := apperrors.EmptyArgError{Subj: "foo"}

	assert.Equal(t, "foo is empty", err.Error())

	assert.True(t, errors.Is(err, apperrors.EmptyArgError{Subj: "foo"}))
	assert.True(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.EmptyArgError{Subj: "foo"}))

	assert.False(t, errors.Is(err, apperrors.EmptyArgError{Subj: "bar"}))
	assert.False(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.EmptyArgError{Subj: "bar"}))

	assert.True(t, errors.As(err, &apperrors.EmptyArgError{}))
	assert.True(t, errors.As(fmt.Errorf("wrap: %w", err), &apperrors.EmptyArgError{}))
}

func TestInvalidArgError(t *testing.T) {
	err := apperrors.InvalidArgError{Subj: "foo", Reason: "theReason"}

	assert.Equal(t, "invalid foo: theReason", err.Error())

	assert.True(t, errors.Is(err, apperrors.InvalidArgError{Subj: "foo", Reason: "theReason"}))
	assert.True(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.InvalidArgError{Subj: "foo", Reason: "theReason"}))

	assert.False(t, errors.Is(err, apperrors.InvalidArgError{Subj: "bar", Reason: "theReason"}))
	assert.False(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.InvalidArgError{Subj: "bar", Reason: "theReason"}))

	assert.False(t, errors.Is(err, apperrors.InvalidArgError{Subj: "foo", Reason: "theOtherReason"}))
	assert.False(t, errors.Is(fmt.Errorf("wrap: %w", err), apperrors.InvalidArgError{Subj: "foo", Reason: "theOtherReason"}))

	assert.True(t, errors.As(err, &apperrors.InvalidArgError{}))
	assert.True(t, errors.As(fmt.Errorf("wrap: %w", err), &apperrors.InvalidArgError{}))
}
