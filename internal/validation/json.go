package validation

import (
	"encoding/json"
	"regexp"
	"sort"

	"github.com/ashep/go-apperrors"
	"github.com/xeipuuv/gojsonschema"
)

// Schema is a record validation schema bound to an index name pattern.
type Schema struct {
	Pattern string
	Schema  json.RawMessage
}

type schemaEntry struct {
	pattern string
	re      *regexp.Regexp
	raw     json.RawMessage
	loader  gojsonschema.JSONLoader
}

type JSONValidator struct {
	entries []schemaEntry
}

func NewJSONValidator(schemas map[string]json.RawMessage) *JSONValidator {
	entries := make([]schemaEntry, 0, len(schemas))
	for pattern, sch := range schemas {
		entries = append(entries, schemaEntry{
			pattern: pattern,
			re:      regexp.MustCompile(pattern),
			raw:     sch,
			loader:  gojsonschema.NewBytesLoader(sch),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].pattern < entries[j].pattern
	})

	return &JSONValidator{
		entries: entries,
	}
}

func (v *JSONValidator) Validate(k, s string) error {
	if s == "" {
		return apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "empty",
		}
	}

	for _, e := range v.entries {
		if !e.re.MatchString(k) {
			continue
		}

		res, err := gojsonschema.Validate(e.loader, gojsonschema.NewBytesLoader([]byte(s)))
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

// SchemasFor returns the schemas whose pattern matches the given index name,
// sorted by pattern. These are exactly the schemas a record pushed to that
// index would be validated against.
func (v *JSONValidator) SchemasFor(name string) []Schema {
	res := make([]Schema, 0, len(v.entries))
	for _, e := range v.entries {
		if e.re.MatchString(name) {
			res = append(res, Schema{Pattern: e.pattern, Schema: e.raw})
		}
	}

	return res
}
