package validation_test

import (
	"testing"

	"github.com/ashep/go-apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/validation"
)

func Test_ValidateJSON(tt *testing.T) {
	tt.Run("NilSchemaMap", func(t *testing.T) {
		v := validation.NewJSONValidator(nil)
		require.NotNil(t, v)
		// Nil schema map should skip validation
		err := v.Validate("test", `{"foo":"bar"}`)
		assert.NoError(t, err)
	})

	tt.Run("EmptySchemaMap", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{})
		require.NotNil(t, v)
		// Empty schema map should skip validation
		err := v.Validate("test", `{"foo":"bar"}`)
		assert.NoError(t, err)
	})

	tt.Run("EmptyJSON", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			".*": `{}`,
		})

		err := v.Validate("test", "")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "empty",
		})
	})

	tt.Run("WhitespaceOnlyJSON", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			".*": `{}`,
		})

		err := v.Validate("test", "   ")
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json schema or data",
			Reason: "EOF",
		})
	})

	tt.Run("MalformedSchema", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			".*": "{]",
		})

		err := v.Validate("test", `{"foo":"bar"}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json schema or data",
			Reason: "invalid character ']' looking for beginning of object key string",
		})
	})

	tt.Run("MalformedData", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			".*": `{}`,
		})

		err := v.Validate("test", `{]`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json schema or data",
			Reason: "invalid character ']' looking for beginning of object key string",
		})
	})

	tt.Run("DataValidationError", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			".*": `{"properties":{"foo":{"type":"number"}}}`,
		})

		err := v.Validate("test", `{"foo":"bar"}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "foo: Invalid type. Expected: number, given: string",
		})
	})

	tt.Run("Ok", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			".*": `{"properties":{"foo":{"type":"string"}}}`,
		})

		err := v.Validate("test", `{"foo":"bar"}`)
		assert.NoError(t, err)
	})

	tt.Run("ValidJSONArray", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^array_.*": `{
				"type": "array",
				"items": {"type": "string"}
			}`,
		})

		err := v.Validate("array_test", `["foo", "bar", "baz"]`)
		assert.NoError(t, err)
	})

	tt.Run("InvalidJSONArrayItems", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^array_.*": `{
				"type": "array",
				"items": {"type": "string"}
			}`,
		})

		err := v.Validate("array_test", `["foo", 123, "baz"]`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "1: Invalid type. Expected: string, given: integer",
		})
	})

	tt.Run("NoMatchingPattern", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^user_.*": `{"properties":{"name":{"type":"string"}}}`,
			"^post_.*": `{"properties":{"title":{"type":"string"}}}`,
		})

		// Key doesn't match any pattern, should skip validation
		err := v.Validate("other_key", `{"invalid":"data"}`)
		assert.NoError(t, err)
	})

	tt.Run("MatchingPatternWithValidData", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^user_.*": `{"properties":{"name":{"type":"string"}}}`,
			"^post_.*": `{"properties":{"title":{"type":"string"}}}`,
		})

		err := v.Validate("user_123", `{"name":"John"}`)
		assert.NoError(t, err)

		err = v.Validate("post_456", `{"title":"My Post"}`)
		assert.NoError(t, err)
	})

	tt.Run("MatchingPatternWithInvalidData", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^user_.*": `{"properties":{"name":{"type":"string"}}}`,
			"^post_.*": `{"properties":{"title":{"type":"string"}}}`,
		})

		err := v.Validate("user_123", `{"name":123}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "name: Invalid type. Expected: string, given: integer",
		})

		err = v.Validate("post_456", `{"title":456}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "title: Invalid type. Expected: string, given: integer",
		})
	})

	tt.Run("ComplexSchemaValidation", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^product_.*": `{
				"type": "object",
				"properties": {
					"id": {"type": "integer"},
					"name": {"type": "string"},
					"price": {"type": "number", "minimum": 0}
				},
				"required": ["id", "name", "price"]
			}`,
		})

		// Valid data
		err := v.Validate("product_123", `{"id":123,"name":"Widget","price":9.99}`)
		assert.NoError(t, err)

		// Missing required field
		err = v.Validate("product_123", `{"id":123,"name":"Widget"}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "(root): price is required",
		})

		// Invalid type
		err = v.Validate("product_123", `{"id":"not-a-number","name":"Widget","price":9.99}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "id: Invalid type. Expected: integer, given: string",
		})

		// Invalid minimum constraint
		err = v.Validate("product_123", `{"id":123,"name":"Widget","price":-5}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "price: Must be greater than or equal to 0",
		})
	})

	tt.Run("MultipleMatchingPatterns", func(t *testing.T) {
		// When multiple patterns match, the first matching one should be used
		// (based on map iteration, which is non-deterministic, but at least one should match)
		v := validation.NewJSONValidator(map[string]string{
			".*":      `{"properties":{"any":{"type":"string"}}}`,
			"^test.*": `{"properties":{"test":{"type":"number"}}}`,
		})

		// This should match one of the patterns and validate accordingly
		// Since map iteration is non-deterministic, we can't predict which pattern wins
		// But the validation should still work with one of them
		err := v.Validate("test_key", `{"test":123}`)
		// This might pass or fail depending on which pattern is checked first
		// If ".*" pattern is checked first and doesn't find "any" field, it passes (schema doesn't require it)
		// If "^test.*" is checked first and validates the "test" field, it passes
		assert.NoError(t, err)
	})

	tt.Run("EmptySchemaString", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			".*": `{}`,
		})

		// Empty schema {} should accept any valid JSON
		err := v.Validate("test", `{"anything":"goes"}`)
		assert.NoError(t, err)

		err = v.Validate("test", `{"numbers":123,"nested":{"deep":true}}`)
		assert.NoError(t, err)
	})

	tt.Run("AdditionalPropertiesValidation", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^strict_.*": `{
				"type": "object",
				"properties": {
					"name": {"type": "string"}
				},
				"additionalProperties": false
			}`,
		})

		// Valid - only defined properties
		err := v.Validate("strict_obj", `{"name":"test"}`)
		assert.NoError(t, err)

		// Invalid - additional properties not allowed
		err = v.Validate("strict_obj", `{"name":"test","extra":"field"}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "(root): Additional property extra is not allowed",
		})
	})

	tt.Run("NullValue", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			".*": `{
				"type": "object",
				"properties": {
					"nullable": {"type": ["string", "null"]}
				}
			}`,
		})

		// Valid null value
		err := v.Validate("test", `{"nullable":null}`)
		assert.NoError(t, err)

		// Valid string value
		err = v.Validate("test", `{"nullable":"value"}`)
		assert.NoError(t, err)
	})

	tt.Run("PrimitiveJSONValues", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^number_.*": `{"type": "number"}`,
			"^string_.*": `{"type": "string"}`,
			"^bool_.*":   `{"type": "boolean"}`,
		})

		// Valid primitive values
		err := v.Validate("number_field", `42`)
		assert.NoError(t, err)

		err = v.Validate("number_field", `3.14`)
		assert.NoError(t, err)

		err = v.Validate("string_field", `"hello"`)
		assert.NoError(t, err)

		err = v.Validate("bool_field", `true`)
		assert.NoError(t, err)

		// Invalid primitive values
		err = v.Validate("number_field", `"not a number"`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "(root): Invalid type. Expected: number, given: string",
		})
	})

	tt.Run("NestedObjectValidation", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^user_.*": `{
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"address": {
						"type": "object",
						"properties": {
							"street": {"type": "string"},
							"city": {"type": "string"}
						},
						"required": ["city"]
					}
				}
			}`,
		})

		// Valid nested object
		err := v.Validate("user_123", `{
			"name": "John",
			"address": {
				"street": "Main St",
				"city": "New York"
			}
		}`)
		assert.NoError(t, err)

		// Missing required nested field
		err = v.Validate("user_123", `{
			"name": "John",
			"address": {
				"street": "Main St"
			}
		}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "address: city is required",
		})
	})

	tt.Run("EnumValidation", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^status_.*": `{
				"type": "object",
				"properties": {
					"status": {
						"type": "string",
						"enum": ["active", "inactive", "pending"]
					}
				}
			}`,
		})

		// Valid enum value
		err := v.Validate("status_check", `{"status":"active"}`)
		assert.NoError(t, err)

		// Invalid enum value
		err = v.Validate("status_check", `{"status":"invalid"}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "status: status must be one of the following: \"active\", \"inactive\", \"pending\"",
		})
	})

	tt.Run("StringFormatValidation", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"^email_.*": `{
				"type": "object",
				"properties": {
					"email": {
						"type": "string",
						"format": "email"
					}
				}
			}`,
		})

		// Valid email format
		err := v.Validate("email_field", `{"email":"user@example.com"}`)
		assert.NoError(t, err)

		// Invalid email format
		err = v.Validate("email_field", `{"email":"not-an-email"}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "email: Does not match format 'email'",
		})
	})

	tt.Run("CaseInsensitivePatternMatching", func(t *testing.T) {
		v := validation.NewJSONValidator(map[string]string{
			"(?i)^user_.*": `{
				"type": "object",
				"properties": {
					"name": {"type": "string"}
				},
				"required": ["name"]
			}`,
		})

		// Should match with case-insensitive pattern
		err := v.Validate("USER_123", `{"name":"John"}`)
		assert.NoError(t, err)

		err = v.Validate("User_456", `{"name":"Jane"}`)
		assert.NoError(t, err)

		// Should fail validation on missing field
		err = v.Validate("user_789", `{}`)
		assert.ErrorIs(t, err, apperrors.InvalidArgError{
			Subj:   "json",
			Reason: "(root): name is required",
		})
	})
}
