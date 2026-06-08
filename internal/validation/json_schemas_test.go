package validation_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ashep/ujds/internal/validation"
)

func TestJSONValidator_SchemasFor(tt *testing.T) {
	tt.Run("Nil", func(t *testing.T) {
		v := validation.NewJSONValidator(nil)
		assert.Empty(t, v.SchemasFor("books"))
	})

	tt.Run("MatchingOnlySortedByPattern", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]json.RawMessage{
			"books.*":  json.RawMessage(`{"type":"object"}`),
			".*":       json.RawMessage(`{}`),
			"^movies$": json.RawMessage(`{"type":"array"}`),
		})

		assert.Equal(t, []validation.Schema{
			{Pattern: ".*", Schema: json.RawMessage(`{}`)},
			{Pattern: "books.*", Schema: json.RawMessage(`{"type":"object"}`)},
		}, v.SchemasFor("books"))
	})

	tt.Run("CatchAllOnlyWhenNoSpecificMatch", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]json.RawMessage{
			"books.*": json.RawMessage(`{"type":"object"}`),
			".*":      json.RawMessage(`{}`),
		})

		assert.Equal(t, []validation.Schema{
			{Pattern: ".*", Schema: json.RawMessage(`{}`)},
		}, v.SchemasFor("movies"))
	})
}
