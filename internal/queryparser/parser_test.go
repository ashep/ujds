package queryparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/queryparser"
)

// select * from record_log where data->>'brand'='Brooks Brothers' and (data->'srp')::int>150;

func TestParse_Basic(tt *testing.T) {
	tt.Run("SingleIdentifierAndOperatorAndIdentifier", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1`, q.String("data", 1))
		assert.Equal(t, []any{123}, q.Args())
	})

	tt.Run("OperatorLeftWhitespace", func(t *testing.T) {
		q, err := queryparser.Parse(`foo =123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1`, q.String("data", 1))
		assert.Equal(t, []any{123}, q.Args())
	})

	tt.Run("OperatorRightWhitespace", func(t *testing.T) {
		q, err := queryparser.Parse(`foo= 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1`, q.String("data", 1))
		assert.Equal(t, []any{123}, q.Args())
	})

	tt.Run("OperatorBothWhitespace", func(t *testing.T) {
		q, err := queryparser.Parse(`foo = 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1`, q.String("data", 1))
		assert.Equal(t, []any{123}, q.Args())
	})

	tt.Run("DottedIdentifier", func(t *testing.T) {
		q, err := queryparser.Parse(`foo.bar.baz = 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo'->'bar'->'baz')::int = $1`, q.String("data", 1))
		assert.Equal(t, []any{123}, q.Args())
	})
}

func TestParse_SyntaxError(tt *testing.T) {
	tt.Run("OperatorExpected1", func(t *testing.T) {
		_, err := queryparser.Parse(`foo`)
		assert.EqualError(t, err, "incomplete expression")
	})

	tt.Run("OperatorExpected2", func(t *testing.T) {
		_, err := queryparser.Parse(`a foo b`)
		assert.EqualError(t, err, "operator expected at position 2: a ")
	})

	tt.Run("OperatorExpected3", func(t *testing.T) {
		_, err := queryparser.Parse(`foo = 123 &&`)
		assert.EqualError(t, err, "incomplete expression")
	})

	tt.Run("SingleIdentifierAndOperator", func(t *testing.T) {
		_, err := queryparser.Parse(`foo=`)
		assert.EqualError(t, err, "incomplete expression")
	})

	tt.Run("IdentifierExpected1", func(t *testing.T) {
		_, err := queryparser.Parse(`"`)
		assert.EqualError(t, err, "identifier expected at position 0: ")
	})

	tt.Run("IdentifierExpected2", func(t *testing.T) {
		_, err := queryparser.Parse(`=`)
		assert.EqualError(t, err, "identifier expected at position 0: ")
	})

	tt.Run("IdentifierExpected3", func(t *testing.T) {
		_, err := queryparser.Parse(`&&`)
		assert.EqualError(t, err, "identifier expected at position 0: ")
	})

	tt.Run("IdentifierExpected4", func(t *testing.T) {
		_, err := queryparser.Parse(`a = 1 && &&`)
		assert.EqualError(t, err, "identifier expected at position 9: a = 1 && ")
	})

	tt.Run("UnknownOperator", func(t *testing.T) {
		_, err := queryparser.Parse(`a >> b`)
		assert.EqualError(t, err, "unknown operator '>>' at position 4: a >>")
	})

	tt.Run("InvalidDottedIdentifier1", func(t *testing.T) {
		_, err := queryparser.Parse(`.foo.bar = 123`)
		assert.EqualError(t, err, "identifier syntax error at position 0: ")
	})

	tt.Run("InvalidDottedIdentifier2", func(t *testing.T) {
		_, err := queryparser.Parse(`foo.bar. = 123`)
		assert.EqualError(t, err, "identifier syntax error at position 8: foo.bar.")
	})

	tt.Run("InvalidDottedIdentifier3", func(t *testing.T) {
		_, err := queryparser.Parse(`foo..bar..baz = 123`)
		assert.EqualError(t, err, "identifier syntax error at position 5: foo..")
	})
}

func TestParse_Literal(tt *testing.T) {
	tt.Run("UnquotedString", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=bar`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::text = '"' || $1 || '"'`, q.String("data", 1))
		assert.Equal(t, []any{"bar"}, q.Args())
	})

	tt.Run("UnquotedStringPrefixedWithInt", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=12bar`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::text = '"' || $1 || '"'`, q.String("data", 1))
		assert.Equal(t, []any{"12bar"}, q.Args())
	})

	tt.Run("UnquotedStringPrefixedWithFloat", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=12.34bar`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::text = '"' || $1 || '"'`, q.String("data", 1))
		assert.Equal(t, []any{"12.34bar"}, q.Args())
	})

	tt.Run("QuotedString", func(t *testing.T) {
		q, err := queryparser.Parse(`foo="bar"`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::text = '"' || $1 || '"'`, q.String("data", 1))
		assert.Equal(t, []any{"bar"}, q.Args())
	})

	tt.Run("QuotedInt", func(t *testing.T) {
		q, err := queryparser.Parse(`foo="123"`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::text = '"' || $1 || '"'`, q.String("data", 1))
		assert.Equal(t, []any{"123"}, q.Args())
	})

	tt.Run("QuotedFloat", func(t *testing.T) {
		q, err := queryparser.Parse(`foo="123.45"`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::text = '"' || $1 || '"'`, q.String("data", 1))
		assert.Equal(t, []any{"123.45"}, q.Args())
	})

	tt.Run("Int", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1`, q.String("data", 1))
		assert.Equal(t, []any{123}, q.Args())
	})

	tt.Run("Float", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=123.45`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::float = $1`, q.String("data", 1))
		assert.Equal(t, []any{123.45}, q.Args())
	})
}

func TestParse_OperatorCompare(tt *testing.T) {
	tt.Run("Eq", func(t *testing.T) {
		q, err := queryparser.Parse(`foo = 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1`, q.String("data", 1))
	})

	tt.Run("EqEq", func(t *testing.T) {
		q, err := queryparser.Parse(`foo == 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1`, q.String("data", 1))
	})

	tt.Run("Neq", func(t *testing.T) {
		q, err := queryparser.Parse(`foo != 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int != $1`, q.String("data", 1))
	})

	tt.Run("Lt", func(t *testing.T) {
		q, err := queryparser.Parse(`foo<123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int < $1`, q.String("data", 1))
	})

	tt.Run("Lte", func(t *testing.T) {
		q, err := queryparser.Parse(`foo <= 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int <= $1`, q.String("data", 1))
	})

	tt.Run("Gt", func(t *testing.T) {
		q, err := queryparser.Parse(`foo > 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int > $1`, q.String("data", 1))
	})

	tt.Run("Gte", func(t *testing.T) {
		q, err := queryparser.Parse(`foo >= 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int >= $1`, q.String("data", 1))
	})
}

func TestParse_OperatorLogical(tt *testing.T) {
	tt.Run("And", func(t *testing.T) {
		q, err := queryparser.Parse(`foo = 123 && bar = 321`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1 AND (data->'bar')::int = $2`, q.String("data", 1))
	})

	tt.Run("Or", func(t *testing.T) {
		q, err := queryparser.Parse(`foo = 123 || bar = 321`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo')::int = $1 OR (data->'bar')::int = $2`, q.String("data", 1))
	})
}
