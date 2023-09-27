package indexrepository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ashep/ujds/internal/indexrepository"
)

func TestNameValidator_Validate(tt *testing.T) {
	tt.Parallel()

	tt.Run("Empty", func(t *testing.T) {
		t.Parallel()
		v := indexrepository.NewNameValidator()
		assert.EqualError(t, v.Validate(""), "must not be empty")
	})

	tt.Run("RegExpNotMatches", func(t *testing.T) {
		t.Parallel()
		v := indexrepository.NewNameValidator()
		assert.EqualError(t, v.Validate("@bc"), "must match the regexp ^[a-zA-Z0-9.-]{1,255}$")
	})

	tt.Run("DotPrefix", func(t *testing.T) {
		t.Parallel()
		v := indexrepository.NewNameValidator()
		assert.EqualError(t, v.Validate(".foo"), "must not start or end with a dot")
	})

	tt.Run("DotSuffix", func(t *testing.T) {
		t.Parallel()
		v := indexrepository.NewNameValidator()
		assert.EqualError(t, v.Validate("foo."), "must not start or end with a dot")
	})

	tt.Run("DashPrefix", func(t *testing.T) {
		t.Parallel()
		v := indexrepository.NewNameValidator()
		assert.EqualError(t, v.Validate("-foo"), "must not start or end with a dash")
	})

	tt.Run("DashSuffix", func(t *testing.T) {
		t.Parallel()
		v := indexrepository.NewNameValidator()
		assert.EqualError(t, v.Validate("foo-"), "must not start or end with a dash")
	})

	tt.Run("DashSuffix", func(t *testing.T) {
		t.Parallel()
		v := indexrepository.NewNameValidator()
		assert.EqualError(t, v.Validate("foo-"), "must not start or end with a dash")
	})

	tt.Run("ConsecutiveDots", func(t *testing.T) {
		t.Parallel()
		v := indexrepository.NewNameValidator()
		assert.EqualError(t, v.Validate("foo..bar"), "must not contain consecutive dots")
	})
}
